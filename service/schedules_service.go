package service

import (
	"context"
	"errors"
	"time"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

var ErrPolicyScheduleNotFound = repository.ErrPolicyScheduleNotFound
var ErrPolicyScheduleSellerNameRequired = errors.New("seller name required")

type PolicySchedulesService interface {
	GetUserSchedules(ctx context.Context, userID uuid.UUID) ([]contracts.ProfilePolicyScheduleResponse, error)
	CreateSchedule(ctx context.Context, userID uuid.UUID, req *contracts.CreateProfilePolicyScheduleRequest) (*contracts.ProfilePolicyScheduleResponse, error)
	DeleteSchedule(ctx context.Context, userID uuid.UUID, profileID int64) error
	PrioritiseSchedule(ctx context.Context, userID uuid.UUID, profileID int64) (bool, time.Time, error)
	GetDueSchedules(ctx context.Context) ([]contracts.ProfilePolicySchedule, error)
	DriveSchedule(ctx context.Context, userID uuid.UUID, profileID int64, state contracts.PolicyScheduleState, by *int64) (contracts.ProfilePolicySchedule, error)
	ProcessSchedule(ctx context.Context, userID uuid.UUID, profileID int64) (contracts.ProfilePolicySchedule, error)
}
