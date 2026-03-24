package service

import (
	"database/sql"

	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/repository"
)

// NewBidsService creates a new bids service
func NewBidsService(db *sql.DB) BidsService {
	bidsRepo := repository.NewBidsRepository(db)
	return newBidsService(bidsRepo)
}

// NewTokensService creates a new tokens service
func NewTokensService(db *sql.DB) TokensService {
	tokensRepo := repository.NewUserTokensRepository(db)
	return newTokensService(tokensRepo)
}

// NewCampaignsService creates a new campaigns service
func NewCampaignsService(db *sql.DB, adsCfg *config.AmazonAdsConfig) CampaignsService {
	tokensRepo := repository.NewUserTokensRepository(db)
	return newCampaignsService(tokensRepo, adsCfg)
}

// NewAuthService creates a new auth service
func NewAuthService(db *sql.DB, adsCfg *config.AmazonAdsConfig) AuthService {
	tokensRepo := repository.NewUserTokensRepository(db)
	return newAuthService(tokensRepo, adsCfg)
}

// NewLWAStateService creates a new LWA state encryption service
func NewLWAStateService(key [32]byte) (LWAStateService, error) {
	return newLWAStateService(key)
}

// NewAttachmentService creates a new attachment service
func NewAttachmentService(db *sql.DB) AttachmentService {
	attachedPoliciesRepo := repository.NewAttachedPoliciesRepository(db)
	return newAttachmentService(attachedPoliciesRepo)
}
