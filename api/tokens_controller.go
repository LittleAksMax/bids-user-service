package api

import (
	"net/http"

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
