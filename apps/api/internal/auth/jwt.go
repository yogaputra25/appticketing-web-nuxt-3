package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token expired")
)

// JWTManager handles creation & validation of HMAC-SHA256 JWTs.
type JWTManager struct {
	secret   []byte
	expHours int
}

func NewJWTManager(secret string, expHours int) *JWTManager {
	return &JWTManager{
		secret:   []byte(secret),
		expHours: expHours,
	}
}

// GenerateToken creates a JWT for a user.
func (j *JWTManager) GenerateToken(userID uint64, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(time.Duration(j.expHours) * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

// ParseToken validates and returns claims.
func (j *JWTManager) ParseToken(tokenStr string) (userID uint64, role string, err error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, "", ErrExpiredToken
		}
		return 0, "", ErrInvalidToken
	}
	if !tok.Valid {
		return 0, "", ErrInvalidToken
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", ErrInvalidToken
	}
	uidF, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", ErrInvalidToken
	}
	roleS, _ := claims["role"].(string)
	return uint64(uidF), roleS, nil
}

// =====================================================
// Middlewares
// =====================================================

// Authenticator extracts and validates JWT from Authorization header,
// then attaches userID & role to the request context.
func (j *JWTManager) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			writeUnauthorized(w, "missing authorization header")
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if raw == "" || raw == header {
			// No "Bearer " prefix or empty after strip
			writeUnauthorized(w, "invalid authorization header format")
			return
		}

		userID, role, err := j.ParseToken(raw)
		if err != nil {
			writeUnauthorized(w, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware — must be used AFTER Authenticator.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(UserRoleKey).(string)
			for _, want := range roles {
				if role == want {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeForbidden(w, fmt.Sprintf("requires role: %v", roles))
		})
	}
}

// RequireAdmin is sugar for RequireRole("admin").
func RequireAdmin(next http.Handler) http.Handler {
	return RequireRole("admin")(next)
}

// =====================================================
// Context getters
// =====================================================
func UserIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value(UserIDKey).(uint64)
	return v, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(UserRoleKey).(string)
	return v, ok
}

// =====================================================
// HTTP helpers
// =====================================================
func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"unauthorized","message":%q}`, msg)))
}

func writeForbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"forbidden","message":%q}`, msg)))
}
