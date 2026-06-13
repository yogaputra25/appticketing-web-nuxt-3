package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	queuePrefix     = "war:queue:"
	tokenPrefix     = "war:token:"
	sessionPrefix   = "war:session:"
	rateLimitPrefix = "war:ratelimit:"
)

type WarQueue struct {
	rdb             *redis.Client
	sessionTTL      time.Duration
	rateLimitPerMin int
}

func NewWarQueue(rdb *redis.Client, sessionTTLMinutes, rateLimitPerMin int) *WarQueue {
	return &WarQueue{
		rdb:             rdb,
		sessionTTL:      time.Duration(sessionTTLMinutes) * time.Minute,
		rateLimitPerMin: rateLimitPerMin,
	}
}

func queueKey(eventID uint64) string {
	return queuePrefix + strconv.FormatUint(eventID, 10)
}

func tokenKey(token string) string {
	return tokenPrefix + token
}

func sessionKey(token string) string {
	return sessionPrefix + token
}

func rateLimitKey(userID, eventID uint64) string {
	return rateLimitPrefix + strconv.FormatUint(userID, 10) + ":" + strconv.FormatUint(eventID, 10)
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// JoinQueue adds user to the queue sorted set (using timestamp as score).
// Returns current position (0-based) and a queue token.
// If user already has an active token, returns existing position.
func (q *WarQueue) JoinQueue(ctx context.Context, userID, eventID uint64) (int64, string, error) {
	key := queueKey(eventID)
	now := float64(time.Now().UnixMicro())

	exists, err := q.rdb.ZScore(ctx, key, strconv.FormatUint(userID, 10)).Result()
	if err == nil {
		rank, err := q.rdb.ZRank(ctx, key, strconv.FormatUint(userID, 10)).Result()
		if err == nil {
			return rank, "", nil
		}
		_ = exists
	}

	member := strconv.FormatUint(userID, 10)
	if err := q.rdb.ZAdd(ctx, key, redis.Z{Score: now, Member: member}).Err(); err != nil {
		return 0, "", fmt.Errorf("zadd: %w", err)
	}

	rank, err := q.rdb.ZRank(ctx, key, member).Result()
	if err != nil {
		return 0, "", fmt.Errorf("zrank: %w", err)
	}

	token := generateToken()
	tk := tokenKey(token)
	if err := q.rdb.Set(ctx, tk, fmt.Sprintf("%d:%d", userID, eventID), q.sessionTTL).Err(); err != nil {
		return 0, "", fmt.Errorf("save token: %w", err)
	}

	return rank, token, nil
}

// GetQueueStatus returns position, total in queue, and whether the user is ready.
// If ready, also returns a booking session token.
// Returns position = -1 if user is not in queue.
func (q *WarQueue) GetQueueStatus(ctx context.Context, userID, eventID uint64) (position, total int64, isReady bool, sessionToken string, err error) {
	key := queueKey(eventID)
	member := strconv.FormatUint(userID, 10)

	rank, err := q.rdb.ZRank(ctx, key, member).Result()
	if err == redis.Nil {
		return -1, 0, false, "", nil
	}
	if err != nil {
		return 0, 0, false, "", fmt.Errorf("zrank: %w", err)
	}

	total, err = q.rdb.ZCard(ctx, key).Result()
	if err != nil {
		return 0, 0, false, "", fmt.Errorf("zcard: %w", err)
	}

	if rank != 0 {
		return rank, total, false, "", nil
	}

		sessionToken, err = q.IssueBookingSession(ctx, userID, eventID)
	if err != nil {
		return rank, total, false, "", fmt.Errorf("issue session: %w", err)
	}

	q.rdb.ZRem(ctx, key, member)

	return 0, total, true, sessionToken, nil
}

func (q *WarQueue) IssueBookingSession(ctx context.Context, userID, eventID uint64) (string, error) {
	token := generateToken()
	sk := sessionKey(token)
	if err := q.rdb.Set(ctx, sk, fmt.Sprintf("%d:%d", userID, eventID), q.sessionTTL).Err(); err != nil {
		return "", fmt.Errorf("set session: %w", err)
	}
	return token, nil
}

// ValidateBookingSession validates a booking session token and returns userID and eventID.
// Returns nil error and the decoded values if valid.
func (q *WarQueue) ValidateBookingSession(ctx context.Context, token string) (userID, eventID uint64, err error) {
	sk := sessionKey(token)
	val, err := q.rdb.Get(ctx, sk).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("invalid or expired session token")
	}
	uid, eid, err := parseUserEvent(val)
	if err != nil {
		return 0, 0, fmt.Errorf("malformed session data")
	}
	return uid, eid, nil
}

// RemoveFromQueue removes a user from the queue.
func (q *WarQueue) RemoveFromQueue(ctx context.Context, userID, eventID uint64) error {
	key := queueKey(eventID)
	member := strconv.FormatUint(userID, 10)
	return q.rdb.ZRem(ctx, key, member).Err()
}

// ProcessQueue advances the queue: finds the first N users (default 1) and issues booking sessions.
// Returns the number of users advanced.
func (q *WarQueue) ProcessQueue(ctx context.Context, eventID uint64, batchSize int64) (int, error) {
	key := queueKey(eventID)

	users, err := q.rdb.ZRange(ctx, key, 0, batchSize-1).Result()
	if err != nil {
		return 0, fmt.Errorf("zrange: %w", err)
	}
	if len(users) == 0 {
		return 0, nil
	}

	advanced := 0
	for _, member := range users {
		userID, err := strconv.ParseUint(member, 10, 64)
		if err != nil {
			continue
		}

		sessionToken, err := q.IssueBookingSession(ctx, userID, eventID)
		if err != nil {
			continue
		}

		q.rdb.ZRem(ctx, key, member)
		_ = sessionToken
		advanced++
	}

	return advanced, nil
}

// CleanupExpired removes queue entries that have been waiting too long.
// Returns number of entries removed.
func (q *WarQueue) CleanupExpired(ctx context.Context, eventID uint64, maxAge time.Duration) (int64, error) {
	key := queueKey(eventID)
	cutoff := float64(time.Now().Add(-maxAge).UnixMicro())

	count, err := q.rdb.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatFloat(cutoff, 'f', 0, 64)).Result()
	if err != nil {
		return 0, fmt.Errorf("zremrangebyscore: %w", err)
	}
	return count, nil
}

// CheckRateLimit checks if the user has exceeded the rate limit for joining the war.
// Returns true if rate limited.
func (q *WarQueue) CheckRateLimit(ctx context.Context, userID, eventID uint64) (bool, error) {
	key := rateLimitKey(userID, eventID)

	pipe := q.rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 60*time.Second)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("rate limit check: %w", err)
	}

	count, err := incr.Result()
	if err != nil {
		return false, fmt.Errorf("incr: %w", err)
	}

	return count > int64(q.rateLimitPerMin), nil
}

// ScanQueueKeys returns all event IDs that have active queues.
func (q *WarQueue) ScanQueueKeys(ctx context.Context) ([]uint64, error) {
	var cursor uint64
	var keys []string
	var err error

	keys, cursor, err = q.rdb.Scan(ctx, cursor, queuePrefix+"*", 100).Result()
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	var ids []uint64
	for _, k := range keys {
		idStr := k[len(queuePrefix):]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func parseUserEvent(val string) (uint64, uint64, error) {
	parts := strings.SplitN(val, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format: %s", val)
	}
	uid, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid user id: %w", err)
	}
	eid, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid event id: %w", err)
	}
	return uid, eid, nil
}
