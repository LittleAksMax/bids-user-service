package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/LittleAksMax/bids-user-service/cache"
	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
)

type attachmentController struct {
	attachmentService service.AttachmentService
	cache             cache.RequestCache
}

func invalidateCampaignCache(ctx context.Context, requestCache cache.RequestCache, profileID int64) error {
	return requestCache.Delete(ctx, campaignsCacheKey(profileID))
}

func (atc *attachmentController) GetAttachedPolicies(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	profileID, err := strconv.ParseInt(r.PathValue(profileIDPath), 10, 64)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "profileID should be an integer",
		})
		return
	}

	attached, err := atc.attachmentService.GetAttachedPoliciesForProfile(r.Context(), userID, profileID)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data:    attached,
	})
}

func (atc *attachmentController) AttachPolicy(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	reqs := requests.GetRequestBody[[]contracts.AttachPolicyRequest](r)
	if reqs == nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	if err := atc.attachmentService.AttachPolicies(r.Context(), userID, *reqs); err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to attach policies",
		})
		return
	}

	invalidated := make(map[int64]struct{}, len(*reqs))
	for _, req := range *reqs {
		if _, seen := invalidated[req.ProfileID]; seen {
			continue
		}
		if err := invalidateCampaignCache(r.Context(), atc.cache, req.ProfileID); err != nil {
			requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
				Success: false,
				Error:   "policies attached but failed to invalidate campaigns cache",
			})
			return
		}
		invalidated[req.ProfileID] = struct{}{}
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "policies attached successfully",
			"count":   len(*reqs),
		},
	})
}

func (atc *attachmentController) DetachPolicy(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	reqs := requests.GetRequestBody[[]contracts.DetachPolicyRequest](r)
	if reqs == nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	if err := atc.attachmentService.DetachPolicies(r.Context(), userID, *reqs); err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to detach policies",
		})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "policies detached successfully",
			"count":   len(*reqs),
		},
	})
}
