package api

import (
	"net/http"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
	"github.com/google/uuid"
)

type tokensController struct {
	tokensService service.TokensService
}

func (tc *tokensController) GetUserTokens(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(uuidSubjectKey).(uuid.UUID)
	if !ok {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{
			Success: false,
			Error:   "invalid user context",
		})
		return
	}

	tokens, err := tc.tokensService.GetUserTokens(r.Context(), userID)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to get tokens",
		})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data:    tokens,
	})
}

func (tc *tokensController) GetAuthenticatedRegions(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(uuidSubjectKey).(uuid.UUID)
	if !ok {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{
			Success: false,
			Error:   "invalid user context",
		})
		return
	}

	tokens, err := tc.tokensService.GetUserTokens(r.Context(), userID)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to get tokens",
		})
		return
	}

	// Pack authenticated regions into array to return to client
	authenticatedRegions := make(contracts.UserAuthenticatedRegionsResponse, 0, 3) // Currently there are only 3 regions
	if tokens.RefreshTokenEU != nil {
		authenticatedRegions = append(authenticatedRegions, "EU")
	}
	if tokens.RefreshTokenUS != nil {
		authenticatedRegions = append(authenticatedRegions, "US")
	}
	if tokens.RefreshTokenFE != nil {
		authenticatedRegions = append(authenticatedRegions, "FE")
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data:    contracts.UserAuthenticatedRegionsResponse(authenticatedRegions),
	})
}
