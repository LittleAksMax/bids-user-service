package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
	"github.com/google/uuid"
)

type attachmentController struct {
	attachmentService service.AttachmentService
}

func (atc *attachmentController) GetAttachedPolicies(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(uuidSubjectKey).(uuid.UUID)
	profileID, err := strconv.ParseInt(r.PathValue(profileIDPath), 10, 64)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "profileID should be an integer",
		})
	}

	attached, err := atc.attachmentService.GetAttachedPoliciesForProfile(r.Context(), userID, profileID)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data:    attached,
	})
}

func (atc *attachmentController) AttachPolicy(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(uuidSubjectKey).(uuid.UUID)

	var req contracts.AttachPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	adGroupIDEmpty := req.AdGroupID == nil || *req.AdGroupID == ""
	campaignIDEmpty := req.CampaignID == nil || *req.CampaignID == ""

	if adGroupIDEmpty && campaignIDEmpty {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "at least one of adgroup_id or campaign_id is required",
		})
		return
	}

	if req.PolicyID == "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "policy_id is required",
		})
		return
	}
	if req.ProfileID == 0 {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "profile_id is required",
		})
		return
	}

	if !adGroupIDEmpty {
		adGroupID := *req.AdGroupID
		campaignID := ""
		if req.CampaignID != nil {
			campaignID = *req.CampaignID
		}
		err := atc.attachmentService.AttachPolicyToAdgroup(r.Context(), userID, adGroupID, req.PolicyID, req.ProfileID, campaignID, req.IsLive)
		if err != nil {
			requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
				Success: false,
				Error:   "failed to attach policy to adgroup",
			})
			return
		}
		requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
			Success: true,
			Data:    map[string]string{"message": "policy attached to adgroup successfully"},
		})
		return
	}

	if !campaignIDEmpty {
		campaignID := *req.CampaignID
		err := atc.attachmentService.AttachPolicyToCampaign(r.Context(), userID, campaignID, req.PolicyID, req.ProfileID, req.IsLive)
		if err != nil {
			requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
				Success: false,
				Error:   "failed to attach policy to campaign",
			})
			return
		}
		requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
			Success: true,
			Data:    map[string]string{"message": "policy attached to all adgroups in campaign successfully"},
		})
		return
	}
}

func (atc *attachmentController) DetachPolicy(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(uuidSubjectKey).(uuid.UUID)

	var req contracts.DetachPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	adGroupIDEmpty := req.AdGroupID == nil || *req.AdGroupID == ""
	campaignIDEmpty := req.CampaignID == nil || *req.CampaignID == ""

	if adGroupIDEmpty && campaignIDEmpty {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "at least one of adgroup_id or campaign_id is required",
		})
		return
	}

	if !adGroupIDEmpty {
		adGroupID := *req.AdGroupID
		if err := atc.attachmentService.DetachPolicyFromAdgroup(r.Context(), userID, adGroupID); err != nil {
			requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
				Success: false,
				Error:   "failed to detach policy from adgroup",
			})
			return
		}
		requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
			Success: true,
			Data:    map[string]string{"message": "policy detached from adgroup successfully"},
		})
		return
	}

	if !campaignIDEmpty {
		campaignID := *req.CampaignID
		if err := atc.attachmentService.DetachPolicyFromCampaign(r.Context(), userID, campaignID); err != nil {
			requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
				Success: false,
				Error:   "failed to detach policy from campaign",
			})
			return
		}
		requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
			Success: true,
			Data:    map[string]string{"message": "policy detached from campaign successfully"},
		})
		return
	}
}
