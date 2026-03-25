package service

import (
	"context"
	"fmt"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

type attachmentService struct {
	attachedPoliciesRepo repository.AttachedPoliciesRepository
}

func newAttachmentService(attachedPoliciesRepo repository.AttachedPoliciesRepository) AttachmentService {
	return &attachmentService{
		attachedPoliciesRepo: attachedPoliciesRepo,
	}
}

func (s *attachmentService) GetAttachedPoliciesForProfile(ctx context.Context, userID uuid.UUID, profileID int64) ([]contracts.AttachedPolicyDTO, error) {
	attachedPolicies, err := s.attachedPoliciesRepo.GetByProfileID(ctx, userID, profileID)
	if err != nil {
		return nil, err
	}

	attachedDTOs := make([]contracts.AttachedPolicyDTO, len(attachedPolicies))
	for i, attached := range attachedPolicies {
		attachedDTOs[i].PolicyID = attached.PolicyID
		attachedDTOs[i].AdGroupID = attached.AdGroupID
		attachedDTOs[i].CampaignID = attached.CampaignID
		attachedDTOs[i].IsLive = attached.IsLive
	}

	return attachedDTOs, nil
}

func (s *attachmentService) AttachPolicies(ctx context.Context, userID uuid.UUID, reqs []contracts.AttachPolicyRequest) error {
	policies := make([]*contracts.AttachedPolicy, 0, len(reqs))
	for _, req := range reqs {
		policies = append(policies, &contracts.AttachedPolicy{
			AdGroupID:  req.AdGroupID,
			PolicyID:   req.PolicyID,
			UserID:     userID,
			ProfileID:  req.ProfileID,
			CampaignID: req.CampaignID,
			IsLive:     req.IsLive,
		})
	}

	if err := s.attachedPoliciesRepo.UpsertBatch(ctx, policies); err != nil {
		return fmt.Errorf("failed to attach policies: %w", err)
	}

	return nil
}

func (s *attachmentService) DetachPolicies(ctx context.Context, userID uuid.UUID, reqs []contracts.DetachPolicyRequest) error {
	if err := s.attachedPoliciesRepo.DeleteBatch(ctx, userID, reqs); err != nil {
		return fmt.Errorf("failed to detach policies: %w", err)
	}

	return nil
}
