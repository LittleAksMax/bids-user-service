package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Bid represents a bid entry
type Bid struct {
	UserID     uuid.UUID
	ProfileID  int64
	CampaignID string
	AdGroupID  string
	PolicyID   string
	FromBid    float64   // The original bid value
	ToBid      float64   // The new bid value
	ChangeDate time.Time // When the bid change was made
	IsLive     bool      // Whether the bid is currently live
}

// BidsRepository defines the interface for bid operations
type BidsRepository interface {
	// Create inserts a new bid (bids are immutable - cannot be updated or deleted)
	Create(ctx context.Context, bid *Bid) error

	// ListWithFilters retrieves bids with optional filters
	ListWithFilters(ctx context.Context, filters *BidFilters) ([]*Bid, error)
}

// BidFilters represents optional filters for querying bids
type BidFilters struct {
	UserID     *uuid.UUID
	ProfileID  *int64
	CampaignID *string
	AdGroupID  *string
	PolicyID   *string
	StartDate  *time.Time
	EndDate    *time.Time
}

// bidsRepository is the concrete implementation
type bidsRepository struct {
	db *sql.DB
}

// NewBidsRepository creates a new bids repository
func NewBidsRepository(db *sql.DB) BidsRepository {
	return &bidsRepository{db: db}
}

func (r *bidsRepository) Create(ctx context.Context, bid *Bid) error {
	query := `
		INSERT INTO bids (user_id, profile_id, campaign_id, adgroup_id, policy_id, from_bid, to_bid, change_date, is_live)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		bid.UserID,
		bid.ProfileID,
		bid.CampaignID,
		bid.AdGroupID,
		bid.PolicyID,
		bid.FromBid,
		bid.ToBid,
		bid.ChangeDate,
		bid.IsLive,
	)

	if err != nil {
		return fmt.Errorf("insert bid: %w", err)
	}

	return nil
}

func (r *bidsRepository) ListWithFilters(ctx context.Context, filters *BidFilters) ([]*Bid, error) {
	query := `
		SELECT user_id, profile_id, campaign_id, adgroup_id, policy_id, from_bid, to_bid, change_date, is_live
		FROM bids
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if filters.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *filters.UserID)
		argPos++
	}

	if filters.ProfileID != nil {
		query += fmt.Sprintf(" AND profile_id = $%d", argPos)
		args = append(args, *filters.ProfileID)
		argPos++
	}

	if filters.CampaignID != nil {
		query += fmt.Sprintf(" AND campaign_id = $%d", argPos)
		args = append(args, *filters.CampaignID)
		argPos++
	}

	if filters.AdGroupID != nil {
		query += fmt.Sprintf(" AND adgroup_id = $%d", argPos)
		args = append(args, *filters.AdGroupID)
		argPos++
	}

	if filters.PolicyID != nil {
		query += fmt.Sprintf(" AND policy_id = $%d", argPos)
		args = append(args, *filters.PolicyID)
		argPos++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND change_date >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND change_date <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	query += " ORDER BY change_date DESC, adgroup_id"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query bids with filters: %w", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return r.scanBids(rows)
}

// scanBids is a helper to scan multiple rows into Bid structs
func (r *bidsRepository) scanBids(rows *sql.Rows) ([]*Bid, error) {
	bids := []*Bid{}

	for rows.Next() {
		bid := &Bid{}
		err := rows.Scan(
			&bid.UserID,
			&bid.ProfileID,
			&bid.CampaignID,
			&bid.AdGroupID,
			&bid.PolicyID,
			&bid.FromBid,
			&bid.ToBid,
			&bid.ChangeDate,
			&bid.IsLive,
		)
		if err != nil {
			return nil, fmt.Errorf("scan bid: %w", err)
		}
		bids = append(bids, bid)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return bids, nil
}
