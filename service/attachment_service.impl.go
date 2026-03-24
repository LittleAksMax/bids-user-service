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

func (s *attachmentService) AttachPolicyToAdgroup(ctx context.Context, userID uuid.UUID, adGroupID string, policyID string, profileID int64, campaignID string, isLive bool) error {
	policy := &contracts.AttachedPolicy{
		AdGroupID:  adGroupID,
		PolicyID:   policyID,
		UserID:     userID,
		ProfileID:  profileID,
		CampaignID: campaignID,
		IsLive:     isLive,
	}
	if err := s.attachedPoliciesRepo.Upsert(ctx, policy); err != nil {
		return fmt.Errorf("failed to attach policy: %w", err)
	}
	return nil
}

func (s *attachmentService) AttachPolicyToCampaign(ctx context.Context, userID uuid.UUID, campaignID string, policyID string, profileID int64, isLive bool) error {
	adgroups, err := s.attachedPoliciesRepo.GetByCampaignID(ctx, userID, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get adgroups for campaign: %w", err)
	}
	for _, ag := range adgroups {
		policy := &contracts.AttachedPolicy{
			AdGroupID:  ag.AdGroupID,
			PolicyID:   policyID,
			UserID:     userID,
			ProfileID:  profileID,
			CampaignID: ag.CampaignID,
			IsLive:     isLive,
		}
		if err := s.attachedPoliciesRepo.Upsert(ctx, policy); err != nil {
			return fmt.Errorf("failed to attach policy to adgroup %s: %w", ag.AdGroupID, err)
		}
	}
	return nil
}

func (s *attachmentService) DetachPolicyFromAdgroup(ctx context.Context, userID uuid.UUID, adGroupID string) error {
	if err := s.attachedPoliciesRepo.Delete(ctx, userID, adGroupID); err != nil {
		return fmt.Errorf("failed to detach policy for adgroup: %w", err)
	}
	return nil
}

func (s *attachmentService) DetachPolicyFromCampaign(ctx context.Context, userID uuid.UUID, campaignID string) error {
	policies, err := s.attachedPoliciesRepo.GetByCampaignID(ctx, userID, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get policies for campaign: %w", err)
	}
	for _, p := range policies {
		if err := s.attachedPoliciesRepo.Delete(ctx, userID, p.AdGroupID); err != nil {
			return fmt.Errorf("failed to detach policy for adgroup %s: %w", p.AdGroupID, err)
		}
	}
	return nil
}
