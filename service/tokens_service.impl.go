package service

import (
	"context"
	"fmt"

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
