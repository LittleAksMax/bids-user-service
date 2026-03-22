package service

import "github.com/LittleAksMax/bids-user-service/contracts"

// LWAStateService handles encryption and decryption of the OAuth state parameter.
type LWAStateService interface {
	// MarshalAndEncrypt serialises and encrypts the redirect state for use as a query parameter.
	MarshalAndEncrypt(state *contracts.RedirectState) (string, error)

	// DecryptAndParse decodes and decrypts the query parameter back into a redirect state.
	DecryptAndParse(encoded string) (*contracts.RedirectState, error)
}
