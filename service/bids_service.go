package service

import (
	"context"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// BidsSearchOptions represents query options for fetching bids
type BidsSearchOptions struct {
	UserID     uuid.UUID
	ProfileID  int64
	CampaignID *string
	AdGroupID  *string
	Days       int
}

// BidsService handles bid-related operations
type BidsService interface {
	// SearchBids retrieves bids based on search options (profile, campaign, adgroup, days)
	SearchBids(ctx context.Context, opts *BidsSearchOptions) ([]*contracts.BidResponse, error)

	// CreateBid creates a new bid for a user
	CreateBid(ctx context.Context, userID uuid.UUID, req *contracts.CreateBidRequest) (*contracts.BidResponse, error)
}
