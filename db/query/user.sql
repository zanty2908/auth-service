-- name: CreateUser :one
INSERT INTO users (
    id, full_name, email, phone, country, password, aud, role
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND aud = $2 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = $1 AND aud = $2
LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 
OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUser :one
UPDATE users
SET 
    updated_at = now(),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    password = COALESCE(sqlc.narg(password), password)
WHERE 
    id = sqlc.arg(id)
RETURNING *;
