package contracts

import (
	"time"

	"github.com/google/uuid"
)

// BidResponse represents a bid in the response
type BidResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	ProfileID  int64     `json:"profile_id"`
	CampaignID string    `json:"campaign_id"`
	AdGroupID  string    `json:"adgroup_id"`
	PolicyID   string    `json:"policy_id"`
	FromBid    float64   `json:"from_bid"`
	ToBid      float64   `json:"to_bid"`
	ChangeDate time.Time `json:"change_date"`
	IsLive     bool      `json:"is_live"`
}

// UserTokensResponse represents user refresh tokens
type UserTokensResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	RefreshTokenEU *string   `json:"refresh_token_eu"`
	RefreshTokenUS *string   `json:"refresh_token_us"`
	RefreshTokenFE *string   `json:"refresh_token_fe"`
}

type AttachedPolicyDTO struct {
	CampaignID string `json:"campaign_id"`
	AdGroupID  string `json:"adgroup_id"`
	PolicyID   string `json:"policy_id"`
	IsLive     bool   `json:"is_live"`
}

type ProfilePolicyScheduleResponse struct {
	ProfileID       int64     `json:"profile_id"`
	DueAt           time.Time `json:"due_at"`
	IntervalMinutes int64     `json:"interval_minutes"`
}
