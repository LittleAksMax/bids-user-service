package api

import (
	"database/sql"
	"net/http"

	"github.com/LittleAksMax/bids-user-service/cache"
	"github.com/go-chi/chi/v5"

	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/health"
	"github.com/LittleAksMax/bids-util/requests"
)

// NewRouter constructs the main API router by wiring middleware and routes defined elsewhere.
func NewRouter(pool *sql.DB, cache *cache.RedisRefreshStore, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	RegisterMiddleware(r)

	requests.ApplyCORS(
		r,
		cfg.AllowedOrigins,
		[]string{"GET", "POST", "PUT", "DELETE"},
		[]string{"Accept", "Authorization", "Content-Type", "X-Auth-Claims", "X-Auth-Ts", "X-Auth-Sig"},
		[]string{"Set-Cookie"},
		true,
		300,
	)

	// Create health checkers map
	healthCheckers := map[string]health.HealthChecker{
		"database": health.NewDBHealthChecker(pool),
		"cache":    health.NewCacheHealthChecker(cache.Client),
	}

	RegisterRoutes(r, healthCheckers)

	return r
}
