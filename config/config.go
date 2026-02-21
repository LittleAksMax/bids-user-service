package config

import (
	"time"

	"github.com/LittleAksMax/bids-user-service/cache"
	"github.com/LittleAksMax/bids-user-service/db"
	"github.com/LittleAksMax/bids-util/env"
)

type AuthConfig struct {
	AccessTokenSecret string
	SharedSecret      string
	MaxSkew           time.Duration
	ClaimsHeader      string
	TimestampHeader   string
	SignatureHeader   string
}

type Config struct {
	DB    *db.PostgresConnectionConfig
	Cache *cache.RedisConnectionConfig
	Auth  *AuthConfig

	Port int

	PasswordPepper string // Add this field for password pepper

	AllowedOrigins []string // CORS allowed origins, read from ALLOWED_ORIGINS (comma-separated)
}

// Load reads environment variables and returns a Config.
// Required: DATABASE_HOST, DATABASE_PORT, DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME, PORT,
// ACCESS_TOKEN_SECRET, REFRESH_TOKEN_SECRET, VALIDATION_API_KEY, REDIS_HOST, REDIS_PORT, REDIS_PASSWORD
func Load() (*Config, error) {
	return &Config{
		Auth: &AuthConfig{
			AccessTokenSecret: env.GetStrFromEnv("ACCESS_TOKEN_SECRET"),
			SharedSecret:      env.GetStrFromEnv("X_AUTH_SIG_SECRET"),
			MaxSkew:           env.ParseDurationEnv("MAX_SKEW"),
			ClaimsHeader:      env.GetStrFromEnv("CLAIMS_HEADER"),
			TimestampHeader:   env.GetStrFromEnv("TIMESTAMP_HEADER"),
			SignatureHeader:   env.GetStrFromEnv("SIGNATURE_HEADER"),
		},
		DB: &db.PostgresConnectionConfig{
			Host:   env.GetStrFromEnv("DATABASE_HOST"),
			Port:   env.ReadPort("DATABASE_PORT"),
			User:   env.GetStrFromEnv("DATABASE_USER"),
			Passwd: env.GetStrFromEnv("DATABASE_PASSWORD"),
			DBName: env.GetStrFromEnv("DATABASE_NAME"),
		},
		Cache: &cache.RedisConnectionConfig{
			Host:     env.GetStrFromEnv("REDIS_HOST"),
			Port:     env.GetIntFromEnv("REDIS_PORT"),
			Password: env.GetStrFromEnv("REDIS_PASSWORD"),
		},
		Port:           env.ReadPort("PORT"),
		AllowedOrigins: env.GetStrListFromEnv("ALLOWED_ORIGINS"),
	}, nil
}
