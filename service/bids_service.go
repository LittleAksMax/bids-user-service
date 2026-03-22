package service

import (
	"context"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// BidsService handles bid-related operations
type BidsService interface {
	// GetBidsForCampaign retrieves all bids for a specific campaign
	GetBidsForCampaign(ctx context.Context, campaignID string) ([]*contracts.BidResponse, error)

	// CreateBid creates a new bid for a user
	CreateBid(ctx context.Context, userID uuid.UUID, req *contracts.BidRequest) (*contracts.BidResponse, error)
}
