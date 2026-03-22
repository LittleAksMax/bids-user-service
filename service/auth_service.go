package service

import (
	"context"

	"github.com/google/uuid"
)

// AuthService handles Amazon LwA authentication
type AuthService interface {
	// ProcessToken processes the callback from Amazon LwA
	ProcessToken(ctx context.Context, userID uuid.UUID, code, region string) error
}
