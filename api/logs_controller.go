package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
)

type logsController struct {
	logsService service.LogsService
}

func logsProfileIDFromRequest(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue(profileIDPath), 10, 64)
}

func writeInvalidLogsProfileID(w http.ResponseWriter, r *http.Request) {
	requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
		Success: false,
		Error:   fmt.Sprintf("invalid profile ID: %s", r.PathValue(profileIDPath)),
	})
}

func logsPageNumFromRequest(r *http.Request) (int, error) {
	pageNumRaw := r.URL.Query().Get("pageNum")
	if pageNumRaw == "" {
		return 1, nil
	}

	pageNum, err := strconv.Atoi(pageNumRaw)
	if err != nil {
		return 0, err
	}
	if pageNum < 1 {
		return 0, fmt.Errorf("pageNum must be positive")
	}

	return pageNum, nil
}

func (lc *logsController) GetUserLogs(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	profileID, err := logsProfileIDFromRequest(r)
	if err != nil {
		writeInvalidLogsProfileID(w, r)
		return
	}

	pageNum, err := logsPageNumFromRequest(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{Success: false, Error: "invalid pageNum query parameter"})
		return
	}

	logs, err := lc.logsService.GetLogs(r.Context(), userID, profileID, pageNum)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrUserLogsPageNotFound) {
			status = http.StatusNotFound
		}
		requests.WriteJSON(w, status, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{Success: true, Data: logs})
}

func (lc *logsController) CreateUserLog(w http.ResponseWriter, r *http.Request) {
	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	profileID, err := logsProfileIDFromRequest(r)
	if err != nil {
		writeInvalidLogsProfileID(w, r)
		return
	}

	req := requests.GetRequestBody[contracts.CreateUserLogRequest](r)
	if req == nil {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{Success: false, Error: "invalid request"})
		return
	}

	logEntry, err := lc.logsService.CreateLog(r.Context(), userID, profileID, req.Log)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrUserLogRequired) {
			status = http.StatusBadRequest
		}
		requests.WriteJSON(w, status, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	requests.WriteJSON(w, http.StatusCreated, requests.APIResponse{
		Success: true,
		Data:    logEntry,
	})
}
