package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	amazonads "github.com/LittleAksMax/amazon-ads-api-sdk-go"
	amazonadsmodels "github.com/LittleAksMax/amazon-ads-api-sdk-go/models"
	"github.com/LittleAksMax/bids-user-service/config"
	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/repository"
	"github.com/google/uuid"
)

const (
	maxConcurrentRequests = 5
	maxRetries            = 3
)

type campaignsService struct {
	tokensRepo repository.UserTokensRepository
	adsCfg     *config.AmazonAdsConfig
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
		HTTPClient: &http.Client{},
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
			Name:     profiles[0].AccountName,
			Profiles: profiles,
		})
	}

	return sellers, nil
}

var listCampaignOptions = amazonadsmodels.ListCampaignsOptions{
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
		return nil, fmt.Errorf("invalid region %s; maybe no registered refresh token", region)
	}

	adsClient, err := s.newAdsClient(refreshToken, region)
	if err != nil {
		return nil, err
	}

	// Paginate through all campaigns and collect the ones we care about
	campaignPaginator := adsClient.CampaignsService.GetCampaigns(profileID, &listCampaignOptions)

	var filteredCampaigns []amazonadsmodels.Campaign
	for campaignPaginator.HasNext() {
		page, err := campaignPaginator.Next(ctx)
		if err != nil {
			return nil, err
		}
		for _, campaign := range page {
			if campaign.AutoCreationSettings == nil || campaign.MarketplaceScope != amazonadsmodels.MarketplaceScopeSingleMarketplace {
				log.Print("Something went wrong in request, no autoCreationSettings for SP Campaign", campaign.CampaignID)
				continue
			}
			if !campaign.AutoCreationSettings.AutoCreateTargets {
				continue
			}
			filteredCampaigns = append(filteredCampaigns, campaign)
		}
	}

	// Step 2: Fan out goroutines per campaign to fetch ad groups (with semaphore)
	type campaignResult struct {
		index    int
		campaign contracts.Campaign
		err      error
	}

	results := make(chan campaignResult, len(filteredCampaigns))
	sem := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup

	for i, campaign := range filteredCampaigns {
		wg.Add(1)
		go func(i int, camp amazonadsmodels.Campaign) {
			defer wg.Done()

			sem <- struct{}{} // acquire
			adGroupPaginator := adsClient.AdGroupsService.GetAdGroups(profileID, &amazonadsmodels.ListAdGroupsOptions{
				AdProductFilter: amazonadsmodels.Filter[amazonadsmodels.AdProduct]{
					Include: []amazonadsmodels.AdProduct{amazonadsmodels.AdProductSP},
				},
				CampaignIDFilter: &amazonadsmodels.Filter[string]{
					Include: []string{camp.CampaignID},
				},
				StateFilter: &amazonadsmodels.Filter[amazonadsmodels.State]{
					Include: []amazonadsmodels.State{amazonadsmodels.StateEnabled},
				},
			})
			<-sem // release

			// Paginate through all ad groups for this campaign
			var adGroupDTOs []contracts.AdGroup
			for adGroupPaginator.HasNext() {
				sem <- struct{}{} // acquire
				page, err := adGroupPaginator.Next(ctx)
				<-sem // release
				if err != nil {
					results <- campaignResult{i, contracts.Campaign{}, err}
					return
				}
				for _, ag := range page {
					if ag.AdProduct != amazonadsmodels.AdProductSP || ag.MarketplaceScope != amazonadsmodels.MarketplaceScopeSingleMarketplace {
						log.Print("Fetched AdGroup:", ag.AdGroupID, "does not pass assertions")
						continue
					}
					adGroupDTOs = append(adGroupDTOs, contracts.AdGroup{
						ID:           ag.AdGroupID,
						Name:         ag.Name,
						DefaultBid:   ag.Bid.DefaultBid,
						CurrencyCode: ag.Bid.CurrencyCode,
					})
				}
			}

			results <- campaignResult{i, contracts.Campaign{
				ID:       camp.CampaignID,
				Name:     camp.Name,
				AdGroups: adGroupDTOs,
			}, nil}
		}(i, campaign)
	}

	wg.Wait()
	close(results)

	campaignDTOs := make([]contracts.Campaign, len(filteredCampaigns))
	for r := range results {
		if r.err != nil {
			return nil, fmt.Errorf("failed to fetch ad groups for campaign index %d: %w", r.index, r.err)
		}
		campaignDTOs[r.index] = r.campaign
	}

	return campaignDTOs, nil
}
