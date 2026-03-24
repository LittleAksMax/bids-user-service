package service

import (
	"context"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// AttachmentService handles attaching/detaching policies to ad groups
type AttachmentService interface {
	GetAttachedPoliciesForProfile(ctx context.Context, userID uuid.UUID, profileID int64) ([]contracts.AttachedPolicyDTO, error)

	// AttachPolicyToAdgroup attaches a policy to a single adgroup
	AttachPolicyToAdgroup(ctx context.Context, userID uuid.UUID, adGroupID string, policyID string, profileID int64, campaignID string, isLive bool) error

	// AttachPolicyToCampaign attaches a policy to all adgroups in a campaign
	AttachPolicyToCampaign(ctx context.Context, userID uuid.UUID, campaignID string, policyID string, profileID int64, isLive bool) error

	// DetachPolicyFromAdgroup removes an attached policy from a single adgroup
	DetachPolicyFromAdgroup(ctx context.Context, userID uuid.UUID, adGroupID string) error

	// DetachPolicyFromCampaign removes an attached policy from all adgroups in a campaign
	DetachPolicyFromCampaign(ctx context.Context, userID uuid.UUID, campaignID string) error
}
