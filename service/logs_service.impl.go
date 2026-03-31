package service

import (
	"context"
	"strings"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

type logsService struct {
	logsRepo repository.LogsRepository
}

const userLogsPageSize = 100

func newLogsService(logsRepo repository.LogsRepository) LogsService {
	return &logsService{
		logsRepo: logsRepo,
	}
}

func (s *logsService) CreateLog(ctx context.Context, userID uuid.UUID, profileID int64, log string) (*contracts.CreatedUserLogResponse, error) {
	log = strings.TrimSpace(log)
	if log == "" {
		return nil, ErrUserLogRequired
	}

	entry, err := s.logsRepo.Create(ctx, userID, profileID, log)
	if err != nil {
		return nil, err
	}

	return &contracts.CreatedUserLogResponse{
		LogID:     entry.LogID,
		ProfileID: entry.ProfileID,
		Log:       entry.Log,
		Timestamp: entry.Timestamp,
	}, nil
}

func (s *logsService) GetLogs(ctx context.Context, userID uuid.UUID, profileID int64, pageNum int) (*contracts.UserLogsPageResponse, error) {
	offset := (pageNum - 1) * userLogsPageSize
	logs, totalCount, err := s.logsRepo.GetByUserIDAndProfileID(ctx, userID, profileID, userLogsPageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + userLogsPageSize - 1) / userLogsPageSize
	}
	if totalPages == 0 {
		if pageNum > 1 {
			return nil, ErrUserLogsPageNotFound
		}

		return &contracts.UserLogsPageResponse{
			Logs:       []contracts.UserLogResponse{},
			TotalPages: 0,
		}, nil
	}
	if pageNum > totalPages {
		return nil, ErrUserLogsPageNotFound
	}

	responseLogs := make([]contracts.UserLogResponse, 0, len(logs))
	for _, entry := range logs {
		responseLogs = append(responseLogs, contracts.UserLogResponse{
			Log:       entry.Log,
			Timestamp: entry.Timestamp,
		})
	}

	return &contracts.UserLogsPageResponse{
		Logs:       responseLogs,
		TotalPages: totalPages,
	}, nil
}
