-- name: GetOTPAuth :one
SELECT * 
FROM otp_authentications 
WHERE phone = $1 AND expires_at > $2 AND deleted_at is null
ORDER BY created_at DESC
LIMIT 1;

-- name: HasOTPAuthValid :one
SELECT 1
FROM otp_authentications 
WHERE phone = $1 AND expires_at > $2 AND deleted_at is null
ORDER BY created_at DESC;

-- name: CreateOTPAuth :one
INSERT INTO otp_authentications(phone, otp, resend_at, expires_at, created_at)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteOTPAuthByPhone :exec
UPDATE otp_authentications
SET deleted_at = now()
WHERE phone = $1 AND expires_at > $2 AND deleted_at is null;

-- name: DeleteOTPAuthByID :exec
UPDATE otp_authentications 
SET deleted_at = now()
WHERE id = $1;
