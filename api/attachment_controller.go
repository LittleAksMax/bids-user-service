package api

import (
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

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "policies attached successfully",
			"count":   len(*reqs),
		},
	})
}

func (atc *attachmentController) DetachPolicy(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(uuidSubjectKey).(uuid.UUID)

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
