package api

import (
	"encoding/json"
	"net/http"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type bidsController struct {
	bidsService service.BidsService
}

func (bc *bidsController) GetBidsForCampaign(w http.ResponseWriter, r *http.Request) {
	campaignID := chi.URLParam(r, "campaignID")
	if campaignID == "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "missing campaign_id",
		})
		return
	}

	bids, err := bc.bidsService.GetBidsForCampaign(r.Context(), campaignID)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to get bids",
		})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data:    bids,
	})
}

func (bc *bidsController) CreateBid(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(uuidSubjectKey)
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{
			Success: false,
			Error:   "invalid user context",
		})
		return
	}

	var req contracts.BidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	bid, err := bc.bidsService.CreateBid(r.Context(), userID, &req)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to create bid",
		})
		return
	}

	requests.WriteJSON(w, http.StatusCreated, requests.APIResponse{
		Success: true,
		Data:    bid,
	})
}
