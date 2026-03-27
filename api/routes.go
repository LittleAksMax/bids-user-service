package api

import (
	"net/http"

	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-util/requests"
	"github.com/LittleAksMax/bids-util/validation"
	"github.com/go-chi/chi/v5"

	"github.com/LittleAksMax/bids-user-service/health"
)

// Health handler implementation that checks all registered services.
// Always returns success=true since the request itself was fulfilled.
// Individual service health is reported in the data field.
func Health(checkers map[string]health.HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statuses := make(map[string]interface{})
		allHealthy := true

		for name, checker := range checkers {
			if err := checker.HealthCheck(r.Context()); err != nil {
				statuses[name] = map[string]interface{}{
					"status": "unhealthy",
					"error":  err.Error(),
				}
				allHealthy = false
			} else {
				statuses[name] = map[string]interface{}{
					"status": "healthy",
				}
			}
		}

		// Determine HTTP status code based on health
		statusCode := http.StatusOK
		if !allHealthy {
			statusCode = http.StatusServiceUnavailable
		}

		requests.WriteJSON(w, statusCode, requests.APIResponse{
			Success: true,
			Data:    statuses,
		})
	}
}

const uuidSubjectKey = "uuidSubject"

// RegisterRoutes registers all endpoint handlers using the controller methods.
func RegisterRoutes(
	r chi.Router,
	bc bidsController,
	tc tokensController,
	cc campaignsController,
	ac authController,
	atc attachmentController,
	psc policySchedulesController,
	authCfg *config.AuthConfig,
	healthCheckers map[string]health.HealthChecker,
) {
	// Health
	r.Get("/health", Health(healthCheckers))

	r.Route("/lwa", func(r chi.Router) {
		// Public auth callback — locked down to Amazon OAuth origins only
		r.Group(func(r chi.Router) {
			requests.ApplyCORS(
				r,
				[]string{
					"https://eu.account.amazon.com",   // EU
					"https://www.amazon.com",          // US
					"https://apac.account.amazon.com", // FE
				},
				[]string{"GET"},
				[]string{"Accept", "Content-Type"},
				nil,
				false,
				300,
			)
			r.Get("/process_token", ac.ProcessToken)
		})
		r.Group(func(r chi.Router) {
			r.Use(
				requests.ValidateAccessToken(
					authCfg.SharedSecret,
					authCfg.AccessTokenSecret,
					authCfg.MaxSkew,
					authCfg.ClaimsHeader,
					authCfg.TimestampHeader,
					authCfg.SignatureHeader,
				),
				requests.EnsureValidSubject(
					authCfg.ClaimsHeader,
					uuidSubjectKey,
				),
			)
			r.Get("/{region}", ac.LWA)
		})
	})

	validationFuncs := []func(interface{}) error{
		validation.ValidateRequiredFields,
		validation.ValidateUUIDs,
		validation.ValidateNonNegativeFields,
	}

	// Authenticated routes
	r.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(
				requests.ValidateAccessToken(
					authCfg.SharedSecret,
					authCfg.AccessTokenSecret,
					authCfg.MaxSkew,
					authCfg.ClaimsHeader,
					authCfg.TimestampHeader,
					authCfg.SignatureHeader,
				),
				requests.EnsureValidSubject(
					authCfg.ClaimsHeader,
					uuidSubjectKey,
				),
			)

			// Attaching policies
			r.Route("/attach", func(r chi.Router) {
				r.Get("/{profileID}", atc.GetAttachedPolicies)
				r.With(requests.ValidateRequest[[]contracts.AttachPolicyRequest](validationFuncs)).Put("/", atc.AttachPolicy)
				r.With(requests.ValidateRequest[[]contracts.DetachPolicyRequest](validationFuncs)).Delete("/", atc.DetachPolicy)
			})

			// Bids endpoints
			r.Route("/bids", func(r chi.Router) {
				r.Get("/{profileID}", bc.GetBidsForProfile)
				r.Get("/{profileID}/{campaignID}", bc.GetBidsForCampaign)
				r.Get("/{profileID}/{campaignID}/{adGroupID}", bc.GetBidsForAdGroup)
			})

			// Tokens endpoints
			r.Get("/tokens", tc.GetUserTokens)

			r.Route("/profiles", func(r chi.Router) {
				// Campaigns endpoint
				r.Get("/", cc.GetProfiles)
				r.Get("/{region}/{profileID}/campaigns", cc.GetCampaigns)
			})

			r.Route("/schedules", func(r chi.Router) {
				r.Get("/", psc.GetUserSchedules)
				r.With(requests.ValidateRequest[contracts.CreateProfilePolicyScheduleRequest](validationFuncs)).Post("/", psc.CreateUserSchedule)
				r.With(requests.ValidateRequest[contracts.DeleteProfilePolicyScheduleRequest](validationFuncs)).Delete("/", psc.DeleteUserSchedule)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(requests.RequireAPIKey(authCfg.ServiceAPIKey, apiKeyHeader))
			r.With(requests.ValidateRequest[contracts.CreateBidRequest](validationFuncs)).Post("/bids", bc.CreateBid)
		})
	})

	r.Route("/internal/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(
				requests.RequireAPIKey(authCfg.ServiceAPIKey, apiKeyHeader),
				InjectUUIDSubjectFromHeader(serviceUserIDHeader, uuidSubjectKey),
			)
			r.Get("/tokens", tc.GetUserTokens)
		})
	})

	r.Route("/internal/schedules", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(requests.RequireAPIKey(authCfg.ServiceAPIKey, apiKeyHeader))
			r.Get("/due", psc.GetDueSchedules)
			r.With(requests.ValidateRequest[contracts.DriveProfilePolicyScheduleRequest](validationFuncs)).Post("/drive", psc.DriveSchedule)
		})
	})
}
