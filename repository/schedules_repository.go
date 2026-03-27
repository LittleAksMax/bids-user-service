package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

var errPolicyScheduleNotFound = errors.New("policy schedule not found")
var ErrPolicyScheduleNotFound = errPolicyScheduleNotFound

type PolicySchedulesRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*contracts.ProfilePolicySchedule, error)
	Create(ctx context.Context, schedule *contracts.ProfilePolicySchedule) error
	Delete(ctx context.Context, userID uuid.UUID, profileID int64) error
	GetDue(ctx context.Context, now time.Time) ([]*contracts.ProfilePolicySchedule, error)
	Drive(ctx context.Context, userID uuid.UUID, profileID int64, now time.Time) (*contracts.ProfilePolicySchedule, error)
}

type policySchedulesRepository struct {
	db *sql.DB
}

func NewPolicySchedulesRepository(db *sql.DB) PolicySchedulesRepository {
	return &policySchedulesRepository{db: db}
}

func (r *policySchedulesRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*contracts.ProfilePolicySchedule, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT user_id, profile_id, due_at, interval_minutes, is_active
		 FROM policy_schedules
		 WHERE user_id = $1 AND is_active = TRUE
		 ORDER BY profile_id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var schedules []*contracts.ProfilePolicySchedule
	for rows.Next() {
		schedule := &contracts.ProfilePolicySchedule{}
		if err := rows.Scan(&schedule.UserID, &schedule.ProfileID, &schedule.DueAt, &schedule.IntervalMinutes, &schedule.IsActive); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *policySchedulesRepository) Create(ctx context.Context, schedule *contracts.ProfilePolicySchedule) error {
	return r.db.QueryRowContext(
		ctx,
		`INSERT INTO policy_schedules (user_id, profile_id, due_at, interval_minutes, is_active)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, profile_id) DO UPDATE
		 SET interval_minutes = EXCLUDED.interval_minutes,
		     is_active = TRUE
		 RETURNING user_id, profile_id, due_at, interval_minutes, is_active`,
		schedule.UserID,
		schedule.ProfileID,
		schedule.DueAt,
		schedule.IntervalMinutes,
		schedule.IsActive,
	).Scan(
		&schedule.UserID,
		&schedule.ProfileID,
		&schedule.DueAt,
		&schedule.IntervalMinutes,
		&schedule.IsActive,
	)
}

func (r *policySchedulesRepository) Delete(ctx context.Context, userID uuid.UUID, profileID int64) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE policy_schedules
		 SET is_active = FALSE
		 WHERE user_id = $1 AND profile_id = $2`,
		userID,
		profileID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errPolicyScheduleNotFound
	}

	return nil
}

func (r *policySchedulesRepository) GetDue(ctx context.Context, now time.Time) ([]*contracts.ProfilePolicySchedule, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT user_id, profile_id, due_at, interval_minutes, is_active
		 FROM policy_schedules
		 WHERE is_active = TRUE AND due_at <= $1
		 ORDER BY due_at, user_id, profile_id`,
		now,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var schedules []*contracts.ProfilePolicySchedule
	for rows.Next() {
		schedule := &contracts.ProfilePolicySchedule{}
		if err := rows.Scan(&schedule.UserID, &schedule.ProfileID, &schedule.DueAt, &schedule.IntervalMinutes, &schedule.IsActive); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *policySchedulesRepository) Drive(ctx context.Context, userID uuid.UUID, profileID int64, now time.Time) (*contracts.ProfilePolicySchedule, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var intervalMinutes int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT interval_minutes
		 FROM policy_schedules
		 WHERE user_id = $1 AND profile_id = $2 AND is_active = TRUE
		 FOR UPDATE`,
		userID,
		profileID,
	).Scan(&intervalMinutes)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errPolicyScheduleNotFound
		}
		return nil, err
	}

	nextDueAt := now.Add(time.Duration(intervalMinutes) * time.Minute)
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE policy_schedules
		 SET due_at = $3
		 WHERE user_id = $1 AND profile_id = $2`,
		userID,
		profileID,
		nextDueAt,
	); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &contracts.ProfilePolicySchedule{
		UserID:          userID,
		ProfileID:       profileID,
		DueAt:           nextDueAt,
		IntervalMinutes: intervalMinutes,
		IsActive:        true,
	}, nil
}
