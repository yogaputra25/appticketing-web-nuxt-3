package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                 string
	AppPort                int
	DatabaseURL            string
	RedisURL               string
	JWTSecret              string
	JWTExpiresHours        int
	BookingTTLMinutes      int
	QueueSessionTTLMinutes int
	WarRateLimitPerMin     int
	PaymentMode            string
}

func Load() (*Config, error) {
	// .env optional
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:                 getEnv("APP_ENV", "development"),
		AppPort:                getEnvInt("APP_PORT", 8080),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://ticketing:ticketing_secret@localhost:5432/ticketing?sslmode=disable"),
		RedisURL:               getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		JWTExpiresHours:        getEnvInt("JWT_EXPIRES_HOURS", 24),
		BookingTTLMinutes:      getEnvInt("BOOKING_TTL_MINUTES", 10),
		QueueSessionTTLMinutes: getEnvInt("QUEUE_SESSION_TTL_MINUTES", 5),
		WarRateLimitPerMin:     getEnvInt("WAR_RATE_LIMIT_PER_MIN", 5),
		PaymentMode:            getEnv("PAYMENT_MODE", "simulation"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
