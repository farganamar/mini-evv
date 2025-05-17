package dto

import "github.com/guregu/null/v5"

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    null.Time `json:"expires_at"`
	IssuedAt     null.Time `json:"issued_at"`
}
