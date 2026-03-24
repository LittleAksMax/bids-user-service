package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// AttachedPolicy represents a policy attached to an ad group

// AttachedPoliciesRepository defines the interface for attached policy operations
type AttachedPoliciesRepository interface {
	// Upsert inserts or updates an attached policy (keyed by adgroup_id)
	Upsert(ctx context.Context, policy *contracts.AttachedPolicy) error

	// Delete removes an attached policy by ad group ID
	Delete(ctx context.Context, userID uuid.UUID, adGroupID string) error

	// GetByAdGroupID retrieves an attached policy by ad group ID
	GetByAdGroupID(ctx context.Context, userID uuid.UUID, adGroupID string) (*contracts.AttachedPolicy, error)

	// GetByProfileID retrieves all attached policies for a profile
	GetByProfileID(ctx context.Context, userID uuid.UUID, profileID int64) ([]*contracts.AttachedPolicy, error)

	// GetByCampaignID retrieves all attached policies for a campaign
	GetByCampaignID(ctx context.Context, userID uuid.UUID, campaignID string) ([]*contracts.AttachedPolicy, error)
}

type attachedPoliciesRepository struct {
	db *sql.DB
}

// NewAttachedPoliciesRepository creates a new attached policies repository
func NewAttachedPoliciesRepository(db *sql.DB) AttachedPoliciesRepository {
	return &attachedPoliciesRepository{db: db}
}

func (r *attachedPoliciesRepository) Upsert(ctx context.Context, policy *contracts.AttachedPolicy) error {
	query := `
		INSERT INTO attached_policies (adgroup_id, policy_id, user_id, profile_id, campaign_id, is_live)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (adgroup_id)
		DO UPDATE SET policy_id = EXCLUDED.policy_id, is_live = EXCLUDED.is_live
	`

	_, err := r.db.ExecContext(ctx, query,
		policy.AdGroupID,
		policy.PolicyID,
		policy.UserID,
		policy.ProfileID,
		policy.CampaignID,
		policy.IsLive,
	)
	if err != nil {
		return fmt.Errorf("upsert attached policy: %w", err)
	}

	return nil
}

func (r *attachedPoliciesRepository) Delete(ctx context.Context, userID uuid.UUID, adGroupID string) error {
	query := `DELETE FROM attached_policies WHERE adgroup_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, adGroupID, userID)
	if err != nil {
		return fmt.Errorf("delete attached policy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attached policy not found for adgroup: %s", adGroupID)
	}

	return nil
}

func (r *attachedPoliciesRepository) GetByAdGroupID(ctx context.Context, userID uuid.UUID, adGroupID string) (*contracts.AttachedPolicy, error) {
	query := `
		SELECT adgroup_id, policy_id, user_id, profile_id, campaign_id, is_live
		FROM attached_policies
		WHERE adgroup_id = $1 AND user_id = $2
	`

	policy := &contracts.AttachedPolicy{}
	err := r.db.QueryRowContext(ctx, query, adGroupID, userID).Scan(
		&policy.AdGroupID,
		&policy.PolicyID,
		&policy.UserID,
		&policy.ProfileID,
		&policy.CampaignID,
		&policy.IsLive,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("attached policy not found for adgroup: %s", adGroupID)
	}
	if err != nil {
		return nil, fmt.Errorf("query attached policy: %w", err)
	}

	return policy, nil
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

func (r *attachedPoliciesRepository) GetByCampaignID(ctx context.Context, userID uuid.UUID, campaignID string) ([]*contracts.AttachedPolicy, error) {
	query := `
		SELECT adgroup_id, policy_id, user_id, profile_id, campaign_id, is_live
		FROM attached_policies
		WHERE campaign_id = $1 AND user_id = $2
		ORDER BY adgroup_id
	`

	rows, err := r.db.QueryContext(ctx, query, campaignID, userID)
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
