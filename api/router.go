package api

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"net/http"
	"net/url"

	"github.com/LittleAksMax/bids-user-service/cache"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/go-chi/chi/v5"

	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/health"
	"github.com/LittleAksMax/bids-util/requests"
)

// NewRouter constructs the main API router by wiring middleware and routes defined elsewhere.
func NewRouter(pool *sql.DB, redisCache *cache.RedisRefreshStore, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	RegisterMiddleware(r)

	requests.ApplyCORS(
		r,
		cfg.AllowedOrigins,
		[]string{"GET", "POST", "PUT", "DELETE"},
		[]string{"Accept", "Authorization", "Content-Type", cfg.Auth.ClaimsHeader, cfg.Auth.TimestampHeader, cfg.Auth.SignatureHeader},
		[]string{"Set-Cookie"},
		true,
		300,
	)

	// Create health checkers map
	healthCheckers := map[string]health.HealthChecker{
		"database": health.NewDBHealthChecker(pool),
		"cache":    health.NewCacheHealthChecker(redisCache.Client),
	}

	// Initialise services independently
	bidsService := service.NewBidsService(pool)
	tokensService := service.NewTokensService(pool)
	campaignsService := service.NewCampaignsService(pool, cfg.Ads)
	authService := service.NewAuthService(pool, cfg.Ads)
	attachmentService := service.NewAttachmentService(pool)

	// Create request cache from existing Redis connection
	requestCache := cache.NewRedisRequestCache(redisCache.Client)

	// Create controllers with their respective services
	bc := bidsController{
		bidsService: bidsService,
	}
	tc := tokensController{
		tokensService: tokensService,
	}
	cc := campaignsController{
		campaignsService: campaignsService,
		cache:            requestCache,
	}

	atc := attachmentController{
		attachmentService: attachmentService,
	}

	redirectURL, err := url.Parse(cfg.Ads.RedirectURI)
	if err != nil {
		log.Fatal("could not parse redirect URI configuration")
	}

	lwaKey := sha256.Sum256([]byte(cfg.Ads.LWASecret))
	lwaStateService, err := service.NewLWAStateService(lwaKey)
	if err != nil {
		log.Fatal("could not initialise LWA state service")
	}

	ac := authController{
		authService:     authService,
		lwaStateService: lwaStateService,
		clientId:        cfg.Ads.ClientID,
		redirectUri:     redirectURL,
	}

	RegisterRoutes(r, bc, tc, cc, ac, atc, cfg.Auth, healthCheckers)

	return r
}
