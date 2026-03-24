package contracts

// ProcessTokenRequest represents the callback from Amazon LwA
type ProcessTokenRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// BidRequest represents a request to create a bid
type BidRequest struct {
	CampaignID string  `json:"campaign_id"`
	AdGroupID  string  `json:"adgroup_id"`
	PolicyID   string  `json:"policy_id"`
	FromBid    float64 `json:"from_bid"`
	ToBid      float64 `json:"to_bid"`
}

// AttachPolicyRequest represents a request to attach a policy to an ad group or campaign
// At least one of AdGroupID or CampaignID must be provided
// If both are provided, AdGroupID takes precedence
// PolicyID, ProfileID, and IsLive are always required
type AttachPolicyRequest struct {
	AdGroupID  *string `json:"adgroup_id"`
	CampaignID *string `json:"campaign_id"`
	PolicyID   string  `json:"policy_id"`
	ProfileID  int64   `json:"profile_id"`
	IsLive     bool    `json:"is_live"`
}

// DetachPolicyRequest represents a request to detach a policy from an ad group or campaign
// At least one of AdGroupID or CampaignID must be provided
type DetachPolicyRequest struct {
	AdGroupID  *string `json:"adgroup_id"`
	CampaignID *string `json:"campaign_id"`
}
