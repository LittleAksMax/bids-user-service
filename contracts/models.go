package contracts

import (
	"time"

	"github.com/google/uuid"
)

// BidRequest represents a request to create a bid
type BidRequest struct {
	CampaignID string  `json:"campaign_id"`
	AdGroupID  string  `json:"adgroup_id"`
	PolicyID   string  `json:"policy_id"`
	FromBid    float64 `json:"from_bid"`
	ToBid      float64 `json:"to_bid"`
}

// BidResponse represents a bid in the response
type BidResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	CampaignID string    `json:"campaign_id"`
	AdGroupID  string    `json:"adgroup_id"`
	PolicyID   string    `json:"policy_id"`
	FromBid    float64   `json:"from_bid"`
	ToBid      float64   `json:"to_bid"`
	ChangeDate time.Time `json:"change_date"`
}

// UserTokensResponse represents user refresh tokens
type UserTokensResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	RefreshTokenEU *string   `json:"refresh_token_eu"`
	RefreshTokenUS *string   `json:"refresh_token_us"`
	RefreshTokenFE *string   `json:"refresh_token_fe"`
}

// ProcessTokenRequest represents the callback from Amazon LwA
type ProcessTokenRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// Campaign represents a campaign with nested structure
type Campaign struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	AdGroups []AdGroup `json:"adgroups"`
}

// AdGroup represents an ad group with nested ads
type AdGroup struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	DefaultBid   float64 `json:"default_bid"`
	CurrencyCode string  `json:"currency_code"`
}

// Seller represents a seller with their profiles
type Seller struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Profiles []RegionProfile `json:"profiles"`
}

// RegionProfile represents a profile from a specific region
type RegionProfile struct {
	ProfileID   int64  `json:"profile_id"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
}
