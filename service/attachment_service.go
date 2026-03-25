package service

import (
	"context"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// AttachmentService handles attaching/detaching policies to ad groups
type AttachmentService interface {
	GetAttachedPoliciesForProfile(ctx context.Context, userID uuid.UUID, profileID int64) ([]contracts.AttachedPolicyDTO, error)

	// AttachPolicies attaches or updates a batch of policies.
	AttachPolicies(ctx context.Context, userID uuid.UUID, reqs []contracts.AttachPolicyRequest) error

	// DetachPolicies removes a batch of attached policies.
	DetachPolicies(ctx context.Context, userID uuid.UUID, reqs []contracts.DetachPolicyRequest) error
}
