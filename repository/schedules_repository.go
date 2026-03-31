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
	Prioritise(ctx context.Context, userID uuid.UUID, profileID int64, dueAt time.Time) (bool, time.Time, error)
	GetDue(ctx context.Context, now time.Time) ([]*contracts.ProfilePolicySchedule, error)
	Drive(ctx context.Context, userID uuid.UUID, profileID int64, state contracts.PolicyScheduleState, by *int64) (*contracts.ProfilePolicySchedule, error)
	Process(ctx context.Context, userID uuid.UUID, profileID int64) (*contracts.ProfilePolicySchedule, error)
}

type policySchedulesRepository struct {
	db *sql.DB
}

func NewPolicySchedulesRepository(db *sql.DB) PolicySchedulesRepository {
	return &policySchedulesRepository{db: db}
}

type profilePolicyScheduleScanner interface {
	Scan(dest ...any) error
}

func nullableStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}

	return value.String
}

func scanProfilePolicySchedule(scanner profilePolicyScheduleScanner, schedule *contracts.ProfilePolicySchedule) error {
	var state string
	var sellerName sql.NullString
	if err := scanner.Scan(
		&schedule.UserID,
		&schedule.ProfileID,
		&schedule.DueAt,
		&schedule.IntervalMinutes,
		&sellerName,
		&schedule.IsActive,
		&state,
	); err != nil {
		return err
	}

	schedule.SellerName = nullableStringValue(sellerName)
	schedule.State = contracts.PolicyScheduleState(state)

	return nil
}

func (r *policySchedulesRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*contracts.ProfilePolicySchedule, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT user_id, profile_id, due_at, interval_minutes, seller_name, is_active, state
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
		if err := scanProfilePolicySchedule(rows, schedule); err != nil {
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
	// We will UPDATE ON CONFLICT, to avoid users conning the site to get faster runs than they should:
	// we update and keep the due_at field the same as it was before
	return scanProfilePolicySchedule(
		r.db.QueryRowContext(
			ctx,
			`INSERT INTO policy_schedules (user_id, profile_id, due_at, interval_minutes, seller_name, is_active)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (user_id, profile_id) DO UPDATE
			 SET interval_minutes = EXCLUDED.interval_minutes,
			     seller_name = EXCLUDED.seller_name,
			     is_active = TRUE
			 RETURNING user_id, profile_id, due_at, interval_minutes, seller_name, is_active, state`,
			schedule.UserID,
			schedule.ProfileID,
			schedule.DueAt,
			schedule.IntervalMinutes,
			schedule.SellerName,
			schedule.IsActive,
		),
		schedule,
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

func (r *policySchedulesRepository) Prioritise(ctx context.Context, userID uuid.UUID, profileID int64, dueAt time.Time) (bool, time.Time, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, time.Time{}, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var state string
	var existingDueAt time.Time
	err = tx.QueryRowContext(
		ctx,
		`SELECT state, due_at
		 FROM policy_schedules
		 WHERE user_id = $1 AND profile_id = $2 AND is_active = TRUE
		 FOR UPDATE`,
		userID,
		profileID,
	).Scan(&state, &existingDueAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, time.Time{}, errPolicyScheduleNotFound
		}
		return false, time.Time{}, err
	}

	if contracts.PolicyScheduleState(state) == contracts.PolicyScheduleStateProcessing {
		return false, existingDueAt, nil
	}

	finalDueAt := existingDueAt
	if existingDueAt.After(dueAt) {
		finalDueAt = dueAt
	}

	if err := tx.QueryRowContext(
		ctx,
		`UPDATE policy_schedules
		 SET due_at = $3
		 WHERE user_id = $1 AND profile_id = $2 AND is_active = TRUE
		 RETURNING due_at`,
		userID,
		profileID,
		finalDueAt,
	).Scan(&finalDueAt); err != nil {
		return false, time.Time{}, err
	}

	if err := tx.Commit(); err != nil {
		return false, time.Time{}, err
	}
	committed = true

	return true, finalDueAt, nil
}

func (r *policySchedulesRepository) GetDue(ctx context.Context, now time.Time) ([]*contracts.ProfilePolicySchedule, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT user_id, profile_id, due_at, interval_minutes, seller_name, is_active, state
		 FROM policy_schedules
		 WHERE is_active = TRUE AND due_at <= $1 AND state IN ('PENDING', 'FAILED', 'SOME ERRORS')
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
		if err := scanProfilePolicySchedule(rows, schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *policySchedulesRepository) Drive(ctx context.Context, userID uuid.UUID, profileID int64, state contracts.PolicyScheduleState, by *int64) (*contracts.ProfilePolicySchedule, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var intervalMinutes int64
	var dueAt time.Time
	var sellerName sql.NullString
	err = tx.QueryRowContext(
		ctx,
		`SELECT interval_minutes, due_at, seller_name
		 FROM policy_schedules
		 WHERE user_id = $1 AND profile_id = $2 AND is_active = TRUE
		 FOR UPDATE`, // Locks row
		userID,
		profileID,
	).Scan(&intervalMinutes, &dueAt, &sellerName)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errPolicyScheduleNotFound
		}
		return nil, err
	}

	// If amount specified, push back by that amount instead
	var nextDueAt time.Time
	if by != nil {
		nextDueAt = dueAt.Add(time.Duration(*by) * time.Minute)
	} else {
		nextDueAt = dueAt.Add(time.Duration(intervalMinutes) * time.Minute)
	}

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE policy_schedules
		 SET due_at = $3,
		     state = $4
		 WHERE user_id = $1 AND profile_id = $2`,
		userID,
		profileID,
		nextDueAt,
		state,
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
		SellerName:      nullableStringValue(sellerName),
		State:           state,
		IsActive:        true,
	}, nil
}

func (r *policySchedulesRepository) Process(ctx context.Context, userID uuid.UUID, profileID int64) (*contracts.ProfilePolicySchedule, error) {
	schedule := &contracts.ProfilePolicySchedule{}
	err := scanProfilePolicySchedule(
		r.db.QueryRowContext(
			ctx,
			`UPDATE policy_schedules
			 SET state = $3
			 WHERE user_id = $1 AND profile_id = $2 AND is_active = TRUE
			 RETURNING user_id, profile_id, due_at, interval_minutes, seller_name, is_active, state`,
			userID,
			profileID,
			contracts.PolicyScheduleStateProcessing,
		),
		schedule,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errPolicyScheduleNotFound
		}
		return nil, err
	}

	return schedule, nil
}
