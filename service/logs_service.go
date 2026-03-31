package service

import (
	"context"
	"errors"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

var ErrUserLogRequired = errors.New("log required")
var ErrUserLogsPageNotFound = errors.New("logs page not found")

type LogsService interface {
	CreateLog(ctx context.Context, userID uuid.UUID, profileID int64, log string) (*contracts.CreatedUserLogResponse, error)
	GetLogs(ctx context.Context, userID uuid.UUID, profileID int64, pageNum int) (*contracts.UserLogsPageResponse, error)
}
