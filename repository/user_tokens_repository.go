package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// UserTokens represents the user's refresh tokens for different regions
type UserTokens struct {
	UserID         uuid.UUID
	RefreshTokenEU *string
	RefreshTokenUS *string
	RefreshTokenFE *string
}

// UserTokensRepository defines the interface for user token operations
type UserTokensRepository interface {
	// GetByUserID retrieves user tokens by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserTokens, error)

	// SetTokenEU sets the EU refresh token for a user (inserts or updates)
	SetTokenEU(ctx context.Context, userID uuid.UUID, token string) error

	// SetTokenUS sets the US refresh token for a user (inserts or updates)
	SetTokenUS(ctx context.Context, userID uuid.UUID, token string) error

	// SetTokenFE sets the FE refresh token for a user (inserts or updates)
	SetTokenFE(ctx context.Context, userID uuid.UUID, token string) error

	// Delete removes a user and their tokens
	Delete(ctx context.Context, userID uuid.UUID) error
}

// userTokensRepository is the concrete implementation
type userTokensRepository struct {
	db *sql.DB
}

// NewUserTokensRepository creates a new user tokens repository
func NewUserTokensRepository(db *sql.DB) UserTokensRepository {
	return &userTokensRepository{db: db}
}

func (r *userTokensRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*UserTokens, error) {
	query := `
		SELECT user_id, refresh_token_eu, refresh_token_us, refresh_token_fe
		FROM user_tokens
		WHERE user_id = $1
	`

	tokens := &UserTokens{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&tokens.UserID,
		&tokens.RefreshTokenEU,
		&tokens.RefreshTokenUS,
		&tokens.RefreshTokenFE,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("query user tokens: %w", err)
	}

	return tokens, nil
}

func (r *userTokensRepository) SetTokenEU(ctx context.Context, userID uuid.UUID, token string) error {
	query := `
		INSERT INTO user_tokens (user_id, refresh_token_eu, refresh_token_us, refresh_token_fe)
		VALUES ($1, $2, NULL, NULL)
		ON CONFLICT (user_id)
		DO UPDATE SET refresh_token_eu = EXCLUDED.refresh_token_eu
	`

	_, err := r.db.ExecContext(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("set EU token: %w", err)
	}

	return nil
}

func (r *userTokensRepository) SetTokenUS(ctx context.Context, userID uuid.UUID, token string) error {
	query := `
		INSERT INTO user_tokens (user_id, refresh_token_eu, refresh_token_us, refresh_token_fe)
		VALUES ($1, NULL, $2, NULL)
		ON CONFLICT (user_id)
		DO UPDATE SET refresh_token_us = EXCLUDED.refresh_token_us
	`

	_, err := r.db.ExecContext(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("set US token: %w", err)
	}

	return nil
}

func (r *userTokensRepository) SetTokenFE(ctx context.Context, userID uuid.UUID, token string) error {
	query := `
		INSERT INTO user_tokens (user_id, refresh_token_eu, refresh_token_us, refresh_token_fe)
		VALUES ($1, NULL, NULL, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET refresh_token_fe = EXCLUDED.refresh_token_fe
	`

	_, err := r.db.ExecContext(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("set FE token: %w", err)
	}

	return nil
}

func (r *userTokensRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_tokens WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", userID)
	}

	return nil
}
