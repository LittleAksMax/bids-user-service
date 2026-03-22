package service

import (
	"context"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/google/uuid"
)

// TokensService handles user token operations
type TokensService interface {
	// GetUserTokens retrieves all tokens for a user
	GetUserTokens(ctx context.Context, userID uuid.UUID) (*contracts.UserTokensResponse, error)
}
