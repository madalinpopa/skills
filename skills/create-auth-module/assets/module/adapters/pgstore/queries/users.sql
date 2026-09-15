-- name: CreateUser :exec
INSERT INTO auth.users (
    user_uuid, name, email, password_hash, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserByUUID :one
SELECT user_uuid, name, email, password_hash, created_at, updated_at
FROM auth.users
WHERE user_uuid = $1;

-- name: GetUserByEmail :one
SELECT user_uuid, name, email, password_hash, created_at, updated_at
FROM auth.users
WHERE email = $1;
