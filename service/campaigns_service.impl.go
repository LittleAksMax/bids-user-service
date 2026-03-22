package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	amazonads "github.com/LittleAksMax/amazon-ads-api-sdk-go"
	amazonadsmodels "github.com/LittleAksMax/amazon-ads-api-sdk-go/models"
	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

type campaignsService struct {
	tokensRepo repository.UserTokensRepository
	adsCfg     *config.AmazonAdsConfig
	httpClient *http.Client
}

func newCampaignsService(tokensRepo repository.UserTokensRepository, adsCfg *config.AmazonAdsConfig) CampaignsService {
	return &campaignsService{
		tokensRepo: tokensRepo,
		adsCfg:     adsCfg,
	}
}

// newAdsClient creates a short-lived Amazon Ads API client for a single request.
func (s *campaignsService) newAdsClient(refreshToken string, region string) (*amazonads.AmazonAdsAPIClient, error) {
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

	client.SetRefreshToken(refreshToken)
	return client, nil
}

type regionToken struct {
	token  string
	region string
}

// groupProfilesBySellerID bins a slice of profiles by their seller ID
// Returns a map where keys are seller IDs and values are slices of profiles
func groupProfilesBySellerID(profiles []amazonadsmodels.Profile) map[string][]amazonadsmodels.Profile {
	grouped := make(map[string][]amazonadsmodels.Profile)

	for _, profile := range profiles {
		sellerID := profile.GetSellerID()
		grouped[sellerID] = append(grouped[sellerID], profile)
	}

	return grouped
}

func (s *campaignsService) GetProfiles(ctx context.Context, userID uuid.UUID) ([]contracts.Seller, error) {
	tokens, err := s.tokensRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	var rts []regionToken
	if tokens.RefreshTokenEU != nil {
		rts = append(rts, regionToken{*tokens.RefreshTokenEU, amazonads.AmazonRegions.Europe})
	}
	if tokens.RefreshTokenUS != nil {
		rts = append(rts, regionToken{*tokens.RefreshTokenUS, amazonads.AmazonRegions.NorthAmerica})
	}
	if tokens.RefreshTokenFE != nil {
		rts = append(rts, regionToken{*tokens.RefreshTokenFE, amazonads.AmazonRegions.FarEast})
	}

	if len(rts) == 0 {
		return nil, errors.New("no registered refresh token found")
	}

	// Accumulate all profiles across regions, keyed by seller ID
	sellerMap := make(map[string][]contracts.RegionProfile)

	client, err := s.newAdsClient("", amazonads.AmazonRegions.Europe)
	if err != nil {
		return nil, err
	}
	for _, rt := range rts {
		client.SetRefreshToken(rt.token)
		if err = client.SetRegion(rt.region); err != nil {
			return nil, err
		}

		profiles, err := client.GetProfiles(ctx, &amazonadsmodels.ListProfilesOptions{
			AccessLevel:              "edit",
			ProfileTypeFilter:        []string{"seller"},
			ValidPaymentMethodFilter: "true",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get profiles for region %s: %w", rt.region, err)
		}

		grouped := groupProfilesBySellerID(profiles)

		for sellerID, sellerProfiles := range grouped {
			for _, p := range sellerProfiles {
				sellerMap[sellerID] = append(sellerMap[sellerID], contracts.RegionProfile{
					ProfileID:   p.ProfileID,
					CountryCode: p.CountryCode,
					Region:      rt.region,
					AccountID:   p.AccountInfo.ID,
					AccountName: p.AccountInfo.Name,
					AccountType: p.AccountInfo.Type,
				})
			}
		}
	}

	// Convert the map into a slice of Seller
	sellers := make([]contracts.Seller, 0, len(sellerMap))
	for sellerID, profiles := range sellerMap {
		sellers = append(sellers, contracts.Seller{
			ID:       sellerID,
			Profiles: profiles,
		})
	}

	return sellers, nil
}

var listCampaignOptions amazonadsmodels.ListCampaignsOptions = amazonadsmodels.ListCampaignsOptions{
	AdProductFilter: amazonadsmodels.Filter[amazonadsmodels.AdProduct]{
		Include: []amazonadsmodels.AdProduct{amazonadsmodels.AdProductSP},
	},
	StateFilter: &amazonadsmodels.Filter[amazonadsmodels.State]{
		Include: []amazonadsmodels.State{amazonadsmodels.StateEnabled},
	},
}

func (s *campaignsService) GetCampaigns(ctx context.Context, userID uuid.UUID, profileID int64, region string) ([]contracts.Campaign, error) {
	region = strings.ToUpper(region)
	tokens, err := s.tokensRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var refreshToken string
	if region == "EU" && tokens.RefreshTokenEU != nil {
		refreshToken = *tokens.RefreshTokenEU
	} else if region == "US" && tokens.RefreshTokenUS != nil {
		refreshToken = *tokens.RefreshTokenUS
	} else if region == "FE" && tokens.RefreshTokenFE != nil {
		refreshToken = *tokens.RefreshTokenFE
	} else {
		return nil, fmt.Errorf("invalid region %s; maybe no registered refresh refreshToken", region)
	}

	adsClient, err := s.newAdsClient(refreshToken, region)
	if err != nil {
		return nil, err
	}

	campaigns, err := adsClient.CampaignsService.GetCampaigns(ctx, profileID, &listCampaignOptions)
	if err != nil {
		return nil, err
	}

	campaignDTOs := make([]contracts.Campaign, 0, len(campaigns))
	for _, campaign := range campaigns {
		campaignDTOs = append(campaignDTOs, contracts.Campaign{
			ID:       campaign.CampaignID,
			Name:     campaign.Name,
			AdGroups: []contracts.AdGroup{},
		})

		// TODO: Get adgroups
	}

	return campaignDTOs, nil
}
