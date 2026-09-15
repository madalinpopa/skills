-- name: CreateRefreshToken :exec
INSERT INTO auth.refresh_tokens (
    token_uuid, user_uuid, token_hash, expires_at, created_at
) VALUES ($1, $2, $3, $4, $5);

-- name: GetRefreshTokenByHash :one
SELECT token_uuid, user_uuid, token_hash, expires_at, created_at, revoked_at
FROM auth.refresh_tokens
WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE auth.refresh_tokens
SET revoked_at = $2
WHERE token_uuid = $1;

-- name: RevokeRefreshTokenIfActive :execrows
-- Zero affected rows means another request already revoked the token, so the
-- caller lost the rotation race.
UPDATE auth.refresh_tokens
SET revoked_at = $2
WHERE token_uuid = $1 AND revoked_at IS NULL;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE auth.refresh_tokens
SET revoked_at = $2
WHERE user_uuid = $1 AND revoked_at IS NULL;
