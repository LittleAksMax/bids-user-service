package contracts

// CreateBidRequest represents a request to create a bid
type CreateBidRequest struct {
	ProfileID  int64   `json:"profile_id" validate:"nonnegative"`
	CampaignID string  `json:"campaign_id" validate:"required"`
	AdGroupID  string  `json:"adgroup_id" validate:"required"`
	PolicyID   string  `json:"policy_id" validate:"required"`
	FromBid    float64 `json:"from_bid" validate:"nonnegative"`
	ToBid      float64 `json:"to_bid" validate:"nonnegative"`
	IsLive     bool    `json:"is_live"`
}

// AttachPolicyRequest represents a single policy attachment in a batch request
type AttachPolicyRequest struct {
	AdGroupID  string `json:"adgroup_id" validate:"required"`
	CampaignID string `json:"campaign_id" validate:"required"`
	PolicyID   string `json:"policy_id" validate:"required"`
	ProfileID  int64  `json:"profile_id" validate:"nonnegative"`
	IsLive     bool   `json:"is_live"`
}

// DetachPolicyRequest represents a single policy detachment in a batch request
type DetachPolicyRequest struct {
	AdGroupID  string `json:"adgroup_id" validate:"required"`
	CampaignID string `json:"campaign_id" validate:"required"`
	ProfileID  int64  `json:"profile_id" validate:"nonnegative"`
}

type ProcessProfilePolicyScheduleRequest struct {
	UserID    string `json:"user_id" validate:"uuid"`
	ProfileID int64  `json:"profile_id" validate:"nonnegative"`
}

type DriveProfilePolicyScheduleRequest struct {
	ProcessProfilePolicyScheduleRequest
	TimeoutMinutes *int64 `json:"timeout" validate:"nonnegative"`
}

type CreateUserLogRequest struct {
	Log string `json:"log" validate:"required"`
}

type CreateProfilePolicyScheduleRequest struct {
	ProfileID       int64  `json:"profile_id" validate:"nonnegative"`
	IntervalMinutes int64  `json:"interval_minutes" validate:"nonnegative"`
	SellerName      string `json:"seller_name" validate:"required"`
}

type DeleteProfilePolicyScheduleRequest struct {
	ProfileID int64 `json:"profile_id" validate:"nonnegative"`
}
