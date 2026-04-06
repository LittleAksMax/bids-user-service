package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LittleAksMax/bids-user-service/cache"
	"github.com/LittleAksMax/bids-util/env"
	"github.com/joho/godotenv"

	"github.com/LittleAksMax/bids-user-service/api"
	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/db"
)

const (
	ModeDevelopment = "development"
	ModeProduction  = "production"
)

func main() {
	ctx := context.Background()

	// Load development override file BEFORE config parsing if MODE indicates development.
	mode := env.GetStrFromEnv("MODE")
	if mode != ModeDevelopment && mode != ModeProduction {
		log.Fatalf("invalid environment variable MODE: %s", mode)
	}
	if mode == ModeDevelopment {
		if err := godotenv.Load(".env.Dev"); err != nil {
			log.Fatalf("Failed to load .env.Dev: %v", err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	pool, err := db.Connect(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer func() {
		if err := pool.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	redisCache, err := cache.NewRedisRefreshStore(ctx, cfg.Cache)
	if err != nil {
		log.Fatalf("cache connect error: %v", err)
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			log.Printf("redis failed to close: %v", err)
		}
	}()

	// Migrate automatically
	if err := db.Migrate(cfg.DB.DSN(), "migrations"); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	r := api.NewRouter(pool, redisCache, cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting server on %s (mode=%s)", addr, mode)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
