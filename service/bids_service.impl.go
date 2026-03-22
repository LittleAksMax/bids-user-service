package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

type bidsService struct {
	bidsRepo repository.BidsRepository
}

func newBidsService(bidsRepo repository.BidsRepository) BidsService {
	return &bidsService{
		bidsRepo: bidsRepo,
	}
}

func (s *bidsService) GetBidsForCampaign(ctx context.Context, campaignID string) ([]*contracts.BidResponse, error) {
	bids, err := s.bidsRepo.GetByCampaignID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bids: %w", err)
	}

	responses := make([]*contracts.BidResponse, len(bids))
	for i, bid := range bids {
		responses[i] = &contracts.BidResponse{
			UserID:     bid.UserID,
			CampaignID: bid.CampaignID,
			AdGroupID:  bid.AdGroupID,
			PolicyID:   bid.PolicyID,
			FromBid:    bid.FromBid,
			ToBid:      bid.ToBid,
			ChangeDate: bid.ChangeDate,
		}
	}

	return responses, nil
}

func (s *bidsService) CreateBid(ctx context.Context, userID uuid.UUID, req *contracts.BidRequest) (*contracts.BidResponse, error) {
	bid := &repository.Bid{
		UserID:     userID,
		CampaignID: req.CampaignID,
		AdGroupID:  req.AdGroupID,
		PolicyID:   req.PolicyID,
		FromBid:    req.FromBid,
		ToBid:      req.ToBid,
		ChangeDate: time.Now(),
	}

	if err := s.bidsRepo.Create(ctx, bid); err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	return &contracts.BidResponse{
		UserID:     bid.UserID,
		CampaignID: bid.CampaignID,
		AdGroupID:  bid.AdGroupID,
		PolicyID:   bid.PolicyID,
		FromBid:    bid.FromBid,
		ToBid:      bid.ToBid,
		ChangeDate: bid.ChangeDate,
	}, nil
}
