package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccessTokenReq struct {
	Name        string       `json:"name"`
	UserID      uuid.UUID    `json:"user_id"`
	ExpiresAt   *time.Time   `json:"expires_at"`
	Permissions *Permissions `json:"permissions"`
}

type CreateAccessTokenResp struct {
	ID          uuid.UUID `json:"id"`
	AccessToken string    `json:"access_token"`
}

type ListAccessTokenResp struct {
	Total        int                `json:"total"`
	AccessTokens []*AccessTokenResp `json:"access_tokens"`
}

type AccessTokenResp struct {
	ID            uuid.UUID    `json:"id"`
	Name          string       `json:"name"`
	UserID        uuid.UUID    `json:"user_id"`
	Permissions   *Permissions `json:"permissions"`
	ExpiresAt     *time.Time   `json:"expires_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	InvalidatedAt *time.Time   `json:"invalidated_at,omitempty"`
}
