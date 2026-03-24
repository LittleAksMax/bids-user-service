package contracts

import "github.com/google/uuid"

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

// AttachedPolicy represents a policy attached to an ad group for a user
type AttachedPolicy struct {
	AdGroupID  string
	PolicyID   string
	UserID     uuid.UUID
	ProfileID  int64
	CampaignID string
	IsLive     bool
}
