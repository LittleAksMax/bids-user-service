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

func (s *bidsService) SearchBids(ctx context.Context, opts *BidsSearchOptions) ([]*contracts.BidResponse, error) {
	startDate := time.Now().UTC().AddDate(0, 0, -opts.Days)

	filters := &repository.BidFilters{
		UserID:    &opts.UserID,
		ProfileID: &opts.ProfileID,
		StartDate: &startDate,
	}

	if opts.CampaignID != nil {
		filters.CampaignID = opts.CampaignID
	}
	if opts.AdGroupID != nil {
		filters.AdGroupID = opts.AdGroupID
	}

	bids, err := s.bidsRepo.ListWithFilters(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search bids: %w", err)
	}

	responses := make([]*contracts.BidResponse, 0, len(bids))
	for _, bid := range bids {
		responses = append(responses, &contracts.BidResponse{
			UserID:     bid.UserID,
			ProfileID:  bid.ProfileID,
			CampaignID: bid.CampaignID,
			AdGroupID:  bid.AdGroupID,
			PolicyID:   bid.PolicyID,
			FromBid:    bid.FromBid,
			ToBid:      bid.ToBid,
			ChangeDate: bid.ChangeDate,
			IsLive:     bid.IsLive,
		})
	}

	return responses, nil
}

func (s *bidsService) CreateBid(ctx context.Context, userID uuid.UUID, req *contracts.CreateBidRequest) (*contracts.BidResponse, error) {
	bid := &repository.Bid{
		UserID:     userID,
		ProfileID:  req.ProfileID,
		CampaignID: req.CampaignID,
		AdGroupID:  req.AdGroupID,
		PolicyID:   req.PolicyID,
		FromBid:    req.FromBid,
		ToBid:      req.ToBid,
		ChangeDate: time.Now().UTC(),
		IsLive:     req.IsLive,
	}

	if err := s.bidsRepo.Create(ctx, bid); err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	return &contracts.BidResponse{
		UserID:     bid.UserID,
		ProfileID:  bid.ProfileID,
		CampaignID: bid.CampaignID,
		AdGroupID:  bid.AdGroupID,
		PolicyID:   bid.PolicyID,
		FromBid:    bid.FromBid,
		ToBid:      bid.ToBid,
		ChangeDate: bid.ChangeDate,
		IsLive:     bid.IsLive,
	}, nil
}
