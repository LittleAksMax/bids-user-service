package service

import (
	"context"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// CampaignsService handles campaign retrieval from Amazon Ads API
type CampaignsService interface {
	// GetProfiles retrieves all Amazon Ads *seller* profiles across regions for a user
	GetProfiles(ctx context.Context, userID uuid.UUID) ([]contracts.Seller, error)

	// GetCampaigns retrieves all campaigns for a profile
	GetCampaigns(ctx context.Context, userID uuid.UUID, profileID int64, region string) ([]contracts.Campaign, error)
}
