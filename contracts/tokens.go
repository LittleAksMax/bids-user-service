package contracts

import "github.com/google/uuid"

type RedirectState struct {
	RedirectURL string    `json:"redirect_url"`
	UserID      uuid.UUID `json:"user_id"`
	Region      string    `json:"region"` // Should be EU, US, FE
}
