package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const defaultDays = 30

type bidsController struct {
	bidsService service.BidsService
}

// parseDays extracts the "days" query parameter, defaulting to 30.
func parseDays(r *http.Request) int {
	daysStr := r.URL.Query().Get("days")
	if daysStr == "" {
		return defaultDays
	}
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		return defaultDays
	}
	return days
}

func (bc *bidsController) GetBidsForProfile(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.ParseInt(chi.URLParam(r, "profileID"), 10, 64)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid profile ID: %s", chi.URLParam(r, "profileID")),
		})
		return
	}

	opts := &service.BidsSearchOptions{
		ProfileID: profileID,
		Days:      parseDays(r),
	}

	bids, err := bc.bidsService.SearchBids(r.Context(), opts)
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

func (bc *bidsController) GetBidsForCampaign(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.ParseInt(chi.URLParam(r, "profileID"), 10, 64)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid profile ID: %s", chi.URLParam(r, "profileID")),
		})
		return
	}

	campaignID := chi.URLParam(r, "campaignID")
	if campaignID == "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "missing campaign ID",
		})
		return
	}

	opts := &service.BidsSearchOptions{
		ProfileID:  profileID,
		CampaignID: &campaignID,
		Days:       parseDays(r),
	}

	bids, err := bc.bidsService.SearchBids(r.Context(), opts)
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

func (bc *bidsController) GetBidsForAdGroup(w http.ResponseWriter, r *http.Request) {
	profileID, err := strconv.ParseInt(chi.URLParam(r, "profileID"), 10, 64)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid profile ID: %s", chi.URLParam(r, "profileID")),
		})
		return
	}

	campaignID := chi.URLParam(r, "campaignID")
	if campaignID == "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "missing campaign ID",
		})
		return
	}

	adGroupID := chi.URLParam(r, "adGroupID")
	if adGroupID == "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "missing ad group ID",
		})
		return
	}

	opts := &service.BidsSearchOptions{
		ProfileID:  profileID,
		CampaignID: &campaignID,
		AdGroupID:  &adGroupID,
		Days:       parseDays(r),
	}

	bids, err := bc.bidsService.SearchBids(r.Context(), opts)
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
	// TODO: implement with API key
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
