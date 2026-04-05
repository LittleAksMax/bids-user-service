package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/LittleAksMax/bids-user-service/contracts"
	"github.com/LittleAksMax/bids-user-service/service"
	"github.com/LittleAksMax/bids-util/requests"
)

type authController struct {
	authService     service.AuthService
	lwaStateService service.LWAStateService
	clientId        string
	redirectUri     *url.URL
}

var regionAuthHosts = map[string]string{
	"EU": "eu.account.amazon.com",
	"US": "www.amazon.com",
	"FE": "apac.account.amazon.com",
}

const redirectURIQueryKey = "redirect_uri"

func (ac *authController) LWA(w http.ResponseWriter, r *http.Request) {
	region := strings.ToUpper(r.PathValue("region"))
	host, ok := regionAuthHosts[region]
	if !ok {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "invalid region. Should be one of 'EU', 'US' or 'FE'",
		})
		return
	}

	userID, err := subjectUUIDFromContext(r)
	if err != nil {
		requests.WriteJSON(w, http.StatusUnauthorized, requests.APIResponse{Success: false, Error: err.Error()})
		return
	}

	redirectTo := r.URL.Query().Get(redirectURIQueryKey)
	if redirectTo == "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "redirect_uri query param is empty",
		})
		return
	}

	state := contracts.RedirectState{
		RedirectURL: redirectTo,
		UserID:      userID,
		Region:      region,
	}

	encodedState, err := ac.lwaStateService.MarshalAndEncrypt(&state)
	if err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to encrypt state",
		})
		return
	}

	authURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/ap/oa",
		RawQuery: url.Values{
			"client_id":     {ac.clientId},
			"scope":         {"advertising::campaign_management"},
			"response_type": {"code"},
			"redirect_uri":  {ac.redirectUri.String()}, // This is for the token to be processed
			"state":         {encodedState},
		}.Encode(),
	}

	requests.WriteJSON(w, http.StatusOK, requests.APIResponse{
		Success: true,
		Data: map[string]string{
			"redirectUrl": authURL.String(),
		},
	})
}

func (ac *authController) ProcessToken(w http.ResponseWriter, r *http.Request) {
	// Check for error response
	queryError := r.URL.Query().Get("error")
	if queryError != "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   queryError,
		})
		return
	}

	// Parse query parameters from Amazon LwA callback
	code := r.URL.Query().Get("code")
	queryState := r.URL.Query().Get("state")
	state, err := ac.lwaStateService.DecryptAndParse(queryState)

	if err != nil || code == "" {
		requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
			Success: false,
			Error:   "missing code or state parameter",
		})
		return
	}

	// Now we get the refresh token using the authorisation code and store it in the database
	if err := ac.authService.ProcessToken(r.Context(), state.UserID, code, state.Region); err != nil {
		requests.WriteJSON(w, http.StatusInternalServerError, requests.APIResponse{
			Success: false,
			Error:   "failed to process token",
		})
		return
	}

	// Redirect back to the frontend page that initiated the OAuth flow
	http.Redirect(w, r, state.RedirectURL, http.StatusFound)
}
