// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: oauth_token.sql

package db

import (
	"context"
	"time"
)

const createOAuthToken = `-- name: CreateOAuthToken :one
INSERT INTO oauth_tokens (
    token_id, user_id, device_id, access_token, refresh_token, access_expires_at, refresh_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING token_id, user_id, device_id, access_token, refresh_token, access_expires_at, refresh_expires_at, created_at, deleted_at
`

type CreateOAuthTokenParams struct {
	TokenID          string    `json:"tokenId"`
	UserID           *string   `json:"userId"`
	DeviceID         *string   `json:"deviceId"`
	AccessToken      string    `json:"accessToken"`
	RefreshToken     string    `json:"refreshToken"`
	AccessExpiresAt  time.Time `json:"accessExpiresAt"`
	RefreshExpiresAt time.Time `json:"refreshExpiresAt"`
}

func (q *Queries) CreateOAuthToken(ctx context.Context, arg *CreateOAuthTokenParams) (*OauthToken, error) {
	row := q.db.QueryRow(ctx, createOAuthToken,
		arg.TokenID,
		arg.UserID,
		arg.DeviceID,
		arg.AccessToken,
		arg.RefreshToken,
		arg.AccessExpiresAt,
		arg.RefreshExpiresAt,
	)
	var i OauthToken
	err := row.Scan(
		&i.TokenID,
		&i.UserID,
		&i.DeviceID,
		&i.AccessToken,
		&i.RefreshToken,
		&i.AccessExpiresAt,
		&i.RefreshExpiresAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const deleteOAuthToken = `-- name: DeleteOAuthToken :exec
UPDATE oauth_tokens 
SET deleted_at = now()
WHERE token_id = $1
`

func (q *Queries) DeleteOAuthToken(ctx context.Context, tokenID string) error {
	_, err := q.db.Exec(ctx, deleteOAuthToken, tokenID)
	return err
}

const getOAuthToken = `-- name: GetOAuthToken :one
SELECT token_id, user_id, device_id, access_token, refresh_token, access_expires_at, refresh_expires_at, created_at, deleted_at 
FROM oauth_tokens 
WHERE token_id = $1 AND deleted_at is null
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetOAuthToken(ctx context.Context, tokenID string) (*OauthToken, error) {
	row := q.db.QueryRow(ctx, getOAuthToken, tokenID)
	var i OauthToken
	err := row.Scan(
		&i.TokenID,
		&i.UserID,
		&i.DeviceID,
		&i.AccessToken,
		&i.RefreshToken,
		&i.AccessExpiresAt,
		&i.RefreshExpiresAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return &i, err
}
