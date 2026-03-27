package service

import (
	"context"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

var ErrPolicyScheduleNotFound = repository.ErrPolicyScheduleNotFound

type PolicySchedulesService interface {
	GetUserSchedules(ctx context.Context, userID uuid.UUID) ([]contracts.ProfilePolicyScheduleResponse, error)
	CreateSchedule(ctx context.Context, userID uuid.UUID, req *contracts.CreateProfilePolicyScheduleRequest) (*contracts.ProfilePolicyScheduleResponse, error)
	DeleteSchedule(ctx context.Context, userID uuid.UUID, profileID int64) error
	GetDueSchedules(ctx context.Context) ([]contracts.ProfilePolicySchedule, error)
	DriveSchedule(ctx context.Context, userID uuid.UUID, profileID int64) (contracts.ProfilePolicySchedule, error)
}
