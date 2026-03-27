package api

import (
	"errors"
	"net/http"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
	"github.com/google/uuid"
)

type policySchedulesController struct {
	policySchedulesService service.PolicySchedulesService
}

func (c policySchedulesController) GetUserSchedules(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	schedules, err := c.policySchedulesService.GetUserSchedules(r.Context(), userID)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{Success: true, Data: schedules})
}

func (c policySchedulesController) CreateUserSchedule(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	req := requests.GetRequestBody[contracts.CreateProfilePolicyScheduleRequest](r)
	if req == nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{Success: false, Error: "invalid request"})
		return
	}

	schedule, err := c.policySchedulesService.CreateSchedule(r.Context(), userID, req)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	requests.WriteJSON(w, http.StatusCreated, requests.APIResponse{Success: true, Data: schedule})
}

func (c policySchedulesController) DeleteUserSchedule(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	req := requests.GetRequestBody[contracts.DeleteProfilePolicyScheduleRequest](r)
	if req == nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{Success: false, Error: "invalid request"})
		return
	}

	if err := c.policySchedulesService.DeleteSchedule(r.Context(), userID, req.ProfileID); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrPolicyScheduleNotFound) {
			status = http.StatusNotFound
		}
		requests.WriteJSON(w, status, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "schedule deleted successfully",
		},
	})
}

func (c policySchedulesController) GetDueSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := c.policySchedulesService.GetDueSchedules(r.Context())
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{Success: true, Data: schedules})
}

func (c policySchedulesController) DriveSchedule(w http.ResponseWriter, r *http.Request) {
	req := requests.GetRequestBody[contracts.DriveProfilePolicyScheduleRequest](r)
	if req == nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{Success: false, Error: "invalid request"})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{Success: false, Error: "invalid user ID"})
		return
	}

	schedule, err := c.policySchedulesService.DriveSchedule(r.Context(), userID, req.ProfileID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrPolicyScheduleNotFound) {
			status = http.StatusNotFound
		}
		requests.WriteJSON(w, status, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{Success: true, Data: schedule})
}

func subjectUUIDFromContext(r *http.Request) (uuid.UUID, error) {
	userID, ok := r.Context().Value(uuidSubjectKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("missing user subject")
	}

	return userID, nil
}
