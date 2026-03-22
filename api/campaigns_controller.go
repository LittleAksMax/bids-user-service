package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LittleAksMax/bids-user-service/cache"
	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
	"github.com/google/uuid"
)

const profilesCacheTTL = 15 * time.Minute

type campaignsController struct {
	campaignsService service.CampaignsService
	cache            cache.RequestCache
}

func profilesCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("profiles:%s", userID)
}

func (cc *campaignsController) GetProfiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(uuidSubjectKey).(uuid.UUID)
	cacheKey := profilesCacheKey(userID)

	// Try cache first
	if cached, _, err := cc.cache.Get(r.Context(), cacheKey); err == nil {
		var groupedProfiles []contracts.Seller
		if err := json.Unmarshal([]byte(cached), &groupedProfiles); err == nil {
			requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
				Success: true,
				Data:    groupedProfiles,
			})
			return
		}
	}

	// Cache miss, so fetch from service
	groupedProfiles, err := cc.campaignsService.GetProfiles(r.Context(), userID)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to get profiles",
		})
		return
	}

	// Store in cache (best-effort, so we still serve the request if this fails)
	if data, err := json.Marshal(groupedProfiles); err == nil {
		_ = cc.cache.Set(r.Context(), cacheKey, string(data), profilesCacheTTL)
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data:    groupedProfiles,
	})
}

const profileIDPath = "profileID"
const regionPath = "region"

func (cc *campaignsController) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(uuidSubjectKey).(uuid.UUID)

	region := strings.ToUpper(r.PathValue(regionPath))
	if region != "EU" && region != "US" && region != "FE" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "invalid region: must be EU, US, or FE",
		})
		return
	}

	profileID, err := strconv.ParseInt(r.PathValue(profileIDPath), 10, 64)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to parse profile ID %s", r.PathValue(profileIDPath)),
		})
		return
	}

	campaigns, err := cc.campaignsService.GetCampaigns(r.Context(), userID, profileID, region)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get campaigns for profile ID %d", profileID),
		})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data:    campaigns,
	})
}
