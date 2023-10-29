-- name: CreateUser :one
INSERT INTO users (
    id, full_name, email, birthday, phone, country, password, address, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = $1 LIMIT 1
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
    birthday = COALESCE(sqlc.narg(birthday), birthday),
    avatar = COALESCE(sqlc.narg(avatar), avatar),
    password = COALESCE(sqlc.narg(password), password),
    address = COALESCE(sqlc.narg(address), address),
    gender = COALESCE(sqlc.narg(gender), gender),
    status = COALESCE(sqlc.narg(status), status)
WHERE 
    id = sqlc.arg(id)
RETURNING *;
