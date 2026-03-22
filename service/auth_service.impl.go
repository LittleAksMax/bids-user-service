package service

import (
	"context"
	"fmt"
	"net/http"

	amazonads "github.com/LittleAksMax/amazon-ads-api-sdk-go"
	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

type authService struct {
	tokensRepo repository.UserTokensRepository
	adsCfg     *config.AmazonAdsConfig
	httpClient *http.Client
}

func newAuthService(tokensRepo repository.UserTokensRepository, adsCfg *config.AmazonAdsConfig) AuthService {
	return &authService{
		tokensRepo: tokensRepo,
		adsCfg:     adsCfg,
		httpClient: &http.Client{},
	}
}

// newAdsClient creates a short-lived Amazon Ads API client for a single request.
func (s *authService) newAdsClient(region string) (*amazonads.AmazonAdsAPIClient, error) {
	authCfg := amazonads.NewAmazonAuthAPIConfig(s.adsCfg.ClientID, s.adsCfg.ClientSecret, s.adsCfg.RedirectURI)
	authClient, err := amazonads.NewAmazonAuthClient(authCfg, region)
	if err != nil {
		return nil, fmt.Errorf("create auth client for region %s: %w", region, err)
	}

	client, err := amazonads.NewAmazonAdsAPIClient(&amazonads.Configuration{
		AuthClient: authClient,
		Region:     region,
		HTTPClient: s.httpClient,
	})
	if err != nil {
		return nil, fmt.Errorf("create ads client for region %s: %w", region, err)
	}

	return client, nil
}

func (s *authService) ProcessToken(ctx context.Context, userID uuid.UUID, code, region string) error {
	client, err := s.newAdsClient(region)
	if err != nil {
		return fmt.Errorf("failed to create ads client: %w", err)
	}

	tok, err := client.ExchangeAuthorisationCode(code)
	if err != nil {
		return fmt.Errorf("failed to exchange authorisation code: %w", err)
	}

	switch region {
	case amazonads.AmazonRegions.Europe:
		err = s.tokensRepo.SetTokenEU(ctx, userID, tok.RefreshToken)
	case amazonads.AmazonRegions.NorthAmerica:
		err = s.tokensRepo.SetTokenUS(ctx, userID, tok.RefreshToken)
	case amazonads.AmazonRegions.FarEast:
		err = s.tokensRepo.SetTokenFE(ctx, userID, tok.RefreshToken)
	default:
		return fmt.Errorf("invalid region: %s", region)
	}

	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}
