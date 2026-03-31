package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserLog struct {
	LogID     uuid.UUID
	ProfileID int64
	Log       string
	Timestamp time.Time
}

type LogsRepository interface {
	Create(ctx context.Context, userID uuid.UUID, profileID int64, log string) (*UserLog, error)
	GetByUserIDAndProfileID(ctx context.Context, userID uuid.UUID, profileID int64, limit int, offset int) ([]UserLog, int, error)
}

type logsRepository struct {
	db *sql.DB
}

func NewLogsRepository(db *sql.DB) LogsRepository {
	return &logsRepository{db: db}
}

func (r *logsRepository) Create(ctx context.Context, userID uuid.UUID, profileID int64, log string) (*UserLog, error) {
	var entry UserLog
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO user_events (profile_id, user_id, log)
		 VALUES ($1, $2, $3)
		 RETURNING log_id, profile_id, log, event_at`,
		profileID,
		userID,
		log,
	).Scan(&entry.LogID, &entry.ProfileID, &entry.Log, &entry.Timestamp)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *logsRepository) GetByUserIDAndProfileID(ctx context.Context, userID uuid.UUID, profileID int64, limit int, offset int) ([]UserLog, int, error) {
	var totalCount int
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
		 FROM user_events
		 WHERE user_id = $1 AND profile_id = $2`,
		userID,
		profileID,
	).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT log, event_at
		 FROM user_events
		 WHERE user_id = $1 AND profile_id = $2
		 ORDER BY event_at DESC, log_id DESC
		 LIMIT $3 OFFSET $4`,
		userID,
		profileID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	logs := make([]UserLog, 0)
	for rows.Next() {
		var entry UserLog
		if err := rows.Scan(&entry.Log, &entry.Timestamp); err != nil {
			return nil, 0, err
		}
		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}
