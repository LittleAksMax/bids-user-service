package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

type tokensService struct {
	tokensRepo repository.UserTokensRepository
}

func newTokensService(tokensRepo repository.UserTokensRepository) TokensService {
	return &tokensService{
		tokensRepo: tokensRepo,
	}
}

func (s *tokensService) GetUserTokens(ctx context.Context, userID uuid.UUID) (*contracts.UserTokensResponse, error) {
	tokens, err := s.tokensRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	response := &contracts.UserTokensResponse{
		UserID:         tokens.UserID,
		RefreshTokenEU: tokens.RefreshTokenEU,
		RefreshTokenUS: tokens.RefreshTokenUS,
		RefreshTokenFE: tokens.RefreshTokenFE,
	}

	return response, nil
}

func (s *tokensService) SetUserToken(ctx context.Context, userID uuid.UUID, req *contracts.SetTokenRequest) error {
	var err error

	region := strings.ToUpper(req.Region)
	switch region {
	case "EU":
		err = s.tokensRepo.SetTokenEU(ctx, userID, req.Token)
	case "US":
		err = s.tokensRepo.SetTokenUS(ctx, userID, req.Token)
	case "FE":
		err = s.tokensRepo.SetTokenFE(ctx, userID, req.Token)
	default:
		return fmt.Errorf("invalid region: %s (must be eu, us, or fe)", req.Region)
	}

	if err != nil {
		return fmt.Errorf("failed to set token for region %s: %w", req.Region, err)
	}

	return nil
}
