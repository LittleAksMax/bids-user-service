package contracts

// CreateBidRequest represents a request to create a bid
type CreateBidRequest struct {
	UserID     string  `json:"user_id" validate:"uuid"`
	ProfileID  int64   `json:"profile_id" validate:"nonnegative"`
	CampaignID string  `json:"campaign_id" validate:"required"`
	AdGroupID  string  `json:"adgroup_id" validate:"required"`
	PolicyID   string  `json:"policy_id" validate:"required"`
	FromBid    float64 `json:"from_bid" validate:"nonnegative"`
	ToBid      float64 `json:"to_bid" validate:"nonnegative"`
}

// AttachPolicyRequest represents a request to attach a policy to an ad group or campaign
// At least one of AdGroupID or CampaignID must be provided
// If both are provided, AdGroupID takes precedence
// PolicyID, ProfileID, and IsLive are always required
type AttachPolicyRequest struct {
	AdGroupID  *string `json:"adgroup_id"`
	CampaignID *string `json:"campaign_id"`
	PolicyID   string  `json:"policy_id" validate:"required"`
	ProfileID  int64   `json:"profile_id" validate:"nonnegative"`
	IsLive     bool    `json:"is_live"`
}

// DetachPolicyRequest represents a request to detach a policy from an ad group or campaign
// At least one of AdGroupID or CampaignID must be provided
type DetachPolicyRequest struct {
	AdGroupID  *string `json:"adgroup_id"`
	CampaignID *string `json:"campaign_id"`
}
