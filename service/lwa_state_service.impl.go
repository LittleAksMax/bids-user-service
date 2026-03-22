package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	"github.com/LittleAksMax/bids-user-service/contracts"
)

type lwaStateService struct {
	aead cipher.AEAD
}

func newLWAStateService(key [32]byte) (LWAStateService, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &lwaStateService{aead: aead}, nil
}

// MarshalAndEncrypt serialises the state to JSON, encrypts it with AES-GCM, and returns URL-safe base64.
func (s *lwaStateService) MarshalAndEncrypt(state *contracts.RedirectState) (string, error) {
	plaintext, err := json.Marshal(state)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal appends the ciphertext (with auth tag) to nonce, giving us nonce || ciphertext || tag
	ciphertext := s.aead.Seal(nonce, nonce, plaintext, nil)

	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// DecryptAndParse decodes URL-safe base64, decrypts with AES-GCM, and unmarshals the state.
func (s *lwaStateService) DecryptAndParse(encoded string) (*contracts.RedirectState, error) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	nonceSize := s.aead.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := s.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	var state contracts.RedirectState
	if err := json.Unmarshal(plaintext, &state); err != nil {
		return nil, err
	}

	return &state, nil
}
