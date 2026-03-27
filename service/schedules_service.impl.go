package service

import (
	"context"
	"time"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

type policySchedulesService struct {
	policySchedulesRepository repository.PolicySchedulesRepository
}

func newPolicySchedulesService(policySchedulesRepository repository.PolicySchedulesRepository) PolicySchedulesService {
	return &policySchedulesService{
		policySchedulesRepository: policySchedulesRepository,
	}
}

func (s *policySchedulesService) GetUserSchedules(ctx context.Context, userID uuid.UUID) ([]contracts.ProfilePolicyScheduleResponse, error) {
	schedules, err := s.policySchedulesRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]contracts.ProfilePolicyScheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		response = append(response, contracts.ProfilePolicyScheduleResponse{
			ProfileID:       schedule.ProfileID,
			DueAt:           schedule.DueAt,
			IntervalMinutes: schedule.IntervalMinutes,
		})
	}

	return response, nil
}

func (s *policySchedulesService) CreateSchedule(ctx context.Context, userID uuid.UUID, req *contracts.CreateProfilePolicyScheduleRequest) (*contracts.ProfilePolicyScheduleResponse, error) {
	schedule := contracts.ProfilePolicySchedule{
		UserID:          userID,
		ProfileID:       req.ProfileID,
		DueAt:           time.Now().UTC().Add(10 * time.Minute),
		IntervalMinutes: req.IntervalMinutes,
		IsActive:        true,
	}

	if err := s.policySchedulesRepository.Create(ctx, &schedule); err != nil {
		return nil, err
	}

	return &contracts.ProfilePolicyScheduleResponse{
		ProfileID:       schedule.ProfileID,
		DueAt:           schedule.DueAt,
		IntervalMinutes: schedule.IntervalMinutes,
	}, nil
}

func (s *policySchedulesService) DeleteSchedule(ctx context.Context, userID uuid.UUID, profileID int64) error {
	return s.policySchedulesRepository.Delete(ctx, userID, profileID)
}

func (s *policySchedulesService) GetDueSchedules(ctx context.Context) ([]contracts.ProfilePolicySchedule, error) {
	schedules, err := s.policySchedulesRepository.GetDue(ctx, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	response := make([]contracts.ProfilePolicySchedule, 0, len(schedules))
	for _, schedule := range schedules {
		response = append(response, *schedule)
	}

	return response, nil
}

func (s *policySchedulesService) DriveSchedule(ctx context.Context, userID uuid.UUID, profileID int64) (contracts.ProfilePolicySchedule, error) {
	schedule, err := s.policySchedulesRepository.Drive(ctx, userID, profileID, time.Now().UTC())
	if err != nil {
		return contracts.ProfilePolicySchedule{}, err
	}

	return *schedule, nil
}
