package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Skip("set TEST_REDIS_URL to run integration tests")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatalf("parse redis url: %v", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("ping redis: %v", err)
	}
	return rdb
}

func cleanupQueue(t *testing.T, rdb *redis.Client) {
	t.Helper()
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, "war:*", 100).Iterator()
	for iter.Next(ctx) {
		rdb.Del(ctx, iter.Val())
	}
}

func TestJoinQueue(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()
	eventID := uint64(1)

	pos, token, err := q.JoinQueue(ctx, 100, eventID)
	if err != nil {
		t.Fatalf("join queue: %v", err)
	}
	if pos != 0 {
		t.Errorf("expected position 0, got %d", pos)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	pos, total, isReady, _, err := q.GetQueueStatus(ctx, 100, eventID)
	if err != nil {
		t.Fatalf("get status: %v", err)
	}
	if isReady {
		t.Error("expected not ready (not first in queue)")
	}
	if pos != 0 {
		t.Errorf("expected position 0, got %d", pos)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
}

func TestJoinQueuePosition(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()
	eventID := uint64(2)

	for _, uid := range []uint64{200, 201, 202} {
		_, _, err := q.JoinQueue(ctx, uid, eventID)
		if err != nil {
			t.Fatalf("join user %d: %v", uid, err)
		}
	}

	pos, total, _, _, err := q.GetQueueStatus(ctx, 202, eventID)
	if err != nil {
		t.Fatalf("get status: %v", err)
	}
	if pos != 2 {
		t.Errorf("expected position 2 for third user, got %d", pos)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
}

func TestDuplicateQueueEntry(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()
	eventID := uint64(3)

	pos1, _, err := q.JoinQueue(ctx, 300, eventID)
	if err != nil {
		t.Fatalf("first join: %v", err)
	}

	pos2, _, err := q.JoinQueue(ctx, 300, eventID)
	if err != nil {
		t.Fatalf("second join: %v", err)
	}

	if pos1 != pos2 {
		t.Errorf("expected same position on duplicate join, got %d vs %d", pos1, pos2)
	}
}

func TestProcessQueue(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()
	eventID := uint64(4)

	for _, uid := range []uint64{400, 401, 402} {
		q.JoinQueue(ctx, uid, eventID)
	}

	advanced, err := q.ProcessQueue(ctx, eventID, 1)
	if err != nil {
		t.Fatalf("process queue: %v", err)
	}
	if advanced != 1 {
		t.Errorf("expected 1 advanced, got %d", advanced)
	}

	pos, _, isReady, sessionToken, err := q.GetQueueStatus(ctx, 400, eventID)
	if err != nil {
		t.Fatalf("get status for advanced user: %v", err)
	}
	if !isReady {
		t.Error("expected user 400 to be ready after advance")
	}
	if sessionToken == "" {
		t.Error("expected non-empty session token")
	}

	pos, _, _, _, err = q.GetQueueStatus(ctx, 401, eventID)
	if err != nil {
		t.Fatalf("get status user 401: %v", err)
	}
	if pos != 0 {
		t.Errorf("expected user 401 to be at position 0 now, got %d", pos)
	}
	_ = pos
}

func TestRateLimit(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 2)
	ctx := context.Background()
	eventID := uint64(5)
	userID := uint64(500)

	limited, err := q.CheckRateLimit(ctx, userID, eventID)
	if err != nil {
		t.Fatalf("check rate limit 1: %v", err)
	}
	if limited {
		t.Error("expected not limited on first call")
	}

	limited, err = q.CheckRateLimit(ctx, userID, eventID)
	if err != nil {
		t.Fatalf("check rate limit 2: %v", err)
	}
	if limited {
		t.Error("expected not limited on second call")
	}

	limited, err = q.CheckRateLimit(ctx, userID, eventID)
	if err != nil {
		t.Fatalf("check rate limit 3: %v", err)
	}
	if !limited {
		t.Error("expected limited after 3 calls (limit=2)")
	}
}

func TestNotInQueue(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()

	pos, total, isReady, _, err := q.GetQueueStatus(ctx, 999, 99)
	if err != nil {
		t.Fatalf("get status for non-queued user: %v", err)
	}
	if pos != -1 {
		t.Errorf("expected position -1 for non-queued user, got %d", pos)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
	if isReady {
		t.Error("expected not ready")
	}
}

func TestValidateSessionToken(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()

	token, err := q.IssueBookingSession(ctx, 600, 6)
	if err != nil {
		t.Fatalf("issue session: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// The issueBookingSession is unexported, so we test via GetQueueStatus
	// after advancing a user through the queue
	q.JoinQueue(ctx, 700, 7)
	q.ProcessQueue(ctx, 7, 1)
	_, _, isReady, sessionToken, err := q.GetQueueStatus(ctx, 700, 7)
	if err != nil {
		t.Fatalf("get status: %v", err)
	}
	if !isReady {
		t.Fatal("expected ready after process queue")
	}
	if sessionToken == "" {
		t.Fatal("expected non-empty session token")
	}
}

func TestCleanupExpired(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()
	eventID := uint64(8)

	q.JoinQueue(ctx, 800, eventID)
	q.JoinQueue(ctx, 801, eventID)

	removed, err := q.CleanupExpired(ctx, eventID, 0)
	if err != nil {
		t.Fatalf("cleanup expired: %v", err)
	}
	if removed != 2 {
		t.Errorf("expected 2 removed (maxAge=0), got %d", removed)
	}

	total, err := rdb.ZCard(ctx, "war:queue:8").Result()
	if err != nil {
		t.Fatalf("zcard: %v", err)
	}
	if total != 0 {
		t.Errorf("expected 0 remaining, got %d", total)
	}
}

func TestRemoveFromQueue(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()
	eventID := uint64(9)

	q.JoinQueue(ctx, 900, eventID)

	if err := q.RemoveFromQueue(ctx, 900, eventID); err != nil {
		t.Fatalf("remove from queue: %v", err)
	}

	pos, _, _, _, err := q.GetQueueStatus(ctx, 900, eventID)
	if err != nil {
		t.Fatalf("get status after remove: %v", err)
	}
	if pos != -1 {
		t.Errorf("expected -1 after removal, got %d", pos)
	}
}

func TestJoinQueueThenNotQueuedAfterAdvance(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()
	eventID := uint64(10)

	q.JoinQueue(ctx, 1000, eventID)
	q.ProcessQueue(ctx, eventID, 1)

	pos, _, _, _, err := q.GetQueueStatus(ctx, 1000, eventID)
	if err != nil {
		t.Fatalf("get status after advance: %v", err)
	}
	if pos != -1 {
		t.Errorf("expected -1 after advance, got %d", pos)
	}
}

func TestSessionTokenTTL(t *testing.T) {
	rdb := setupTestRedis(t)
	defer cleanupQueue(t, rdb)

	q := NewWarQueue(rdb, 5, 5)
	ctx := context.Background()

	token, err := q.IssueBookingSession(ctx, 1100, 11)
	if err != nil {
		t.Fatalf("issue session: %v", err)
	}

	ttl, err := rdb.TTL(ctx, "war:session:"+token).Result()
	if err != nil {
		t.Fatalf("get ttl: %v", err)
	}
	if ttl <= 0 || ttl > 6*time.Minute {
		t.Errorf("expected TTL ~5 minutes, got %v", ttl)
	}
}
