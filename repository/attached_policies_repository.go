package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// AttachedPoliciesRepository defines the interface for attached policy operations
type AttachedPoliciesRepository interface {
	// UpsertBatch inserts or updates attached policies atomically.
	UpsertBatch(ctx context.Context, policies []*contracts.AttachedPolicy) error

	// DeleteBatch removes attached policies atomically.
	DeleteBatch(ctx context.Context, userID uuid.UUID, reqs []contracts.DetachPolicyRequest) error

	// GetByProfileID retrieves all attached policies for a profile
	GetByProfileID(ctx context.Context, userID uuid.UUID, profileID int64) ([]*contracts.AttachedPolicy, error)
}

type attachedPoliciesRepository struct {
	db *sql.DB
}

// NewAttachedPoliciesRepository creates a new attached policies repository
func NewAttachedPoliciesRepository(db *sql.DB) AttachedPoliciesRepository {
	return &attachedPoliciesRepository{db: db}
}

func (r *attachedPoliciesRepository) UpsertBatch(ctx context.Context, policies []*contracts.AttachedPolicy) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin upsert attached policies transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	upsertQuery := `
		INSERT INTO attached_policies (adgroup_id, policy_id, user_id, profile_id, campaign_id, is_live)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (adgroup_id)
		DO UPDATE SET
			policy_id = EXCLUDED.policy_id,
			user_id = EXCLUDED.user_id,
			profile_id = EXCLUDED.profile_id,
			campaign_id = EXCLUDED.campaign_id,
			is_live = EXCLUDED.is_live
	`

	for _, policy := range policies {
		if _, err := tx.ExecContext(ctx, upsertQuery,
			policy.AdGroupID,
			policy.PolicyID,
			policy.UserID,
			policy.ProfileID,
			policy.CampaignID,
			policy.IsLive,
		); err != nil {
			return fmt.Errorf("upsert attached policy for adgroup %s: %w", policy.AdGroupID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit upsert attached policies transaction: %w", err)
	}

	return nil
}

func (r *attachedPoliciesRepository) DeleteBatch(ctx context.Context, userID uuid.UUID, reqs []contracts.DetachPolicyRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin delete attached policies transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	deleteQuery := `
		DELETE FROM attached_policies
		WHERE adgroup_id = $1 AND campaign_id = $2 AND profile_id = $3 AND user_id = $4
	`

	for _, req := range reqs {
		result, err := tx.ExecContext(ctx, deleteQuery, req.AdGroupID, req.CampaignID, req.ProfileID, userID)
		if err != nil {
			return fmt.Errorf("delete attached policy for adgroup %s: %w", req.AdGroupID, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("get rows affected for adgroup %s: %w", req.AdGroupID, err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("attached policy not found for adgroup %s", req.AdGroupID)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit delete attached policies transaction: %w", err)
	}

	return nil
}

func (r *attachedPoliciesRepository) GetByProfileID(ctx context.Context, userID uuid.UUID, profileID int64) ([]*contracts.AttachedPolicy, error) {
	query := `
		SELECT adgroup_id, policy_id, user_id, profile_id, campaign_id, is_live
		FROM attached_policies
		WHERE profile_id = $1 AND user_id = $2
		ORDER BY campaign_id, adgroup_id
	`

	rows, err := r.db.QueryContext(ctx, query, profileID, userID)
	if err != nil {
		return nil, fmt.Errorf("query attached policies: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	return r.scanPolicies(rows)
}

func (r *attachedPoliciesRepository) scanPolicies(rows *sql.Rows) ([]*contracts.AttachedPolicy, error) {
	var policies []*contracts.AttachedPolicy

	for rows.Next() {
		policy := &contracts.AttachedPolicy{}
		err := rows.Scan(
			&policy.AdGroupID,
			&policy.PolicyID,
			&policy.UserID,
			&policy.ProfileID,
			&policy.CampaignID,
			&policy.IsLive,
		)
		if err != nil {
			return nil, fmt.Errorf("scan attached policy: %w", err)
		}
		policies = append(policies, policy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return policies, nil
}
