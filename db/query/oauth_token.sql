-- name: CreateOAuthToken :one
INSERT INTO oauth_tokens (
    token_id, aud, platform, user_id, device_id, access_token, refresh_token, access_expires_at, refresh_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetOAuthToken :one
SELECT * 
FROM oauth_tokens 
WHERE token_id = $1 AND deleted_at is null
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteOAuthToken :exec
UPDATE oauth_tokens 
SET deleted_at = now()
WHERE token_id = $1;
