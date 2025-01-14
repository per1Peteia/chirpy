-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
    VALUES ($1, NOW(), NOW(), $2, $3, NULL)
RETURNING
    *;

-- name: GetRefreshToken :one
SELECT
    *
FROM
    refresh_tokens
WHERE
    token = $1
    AND expires_at > now()
    AND revoked_at IS NULL;

-- name: RevokeRefreshToken :one
UPDATE
    refresh_tokens
SET
    updated_at = now(),
    revoked_at = now()
WHERE
    token = $1
RETURNING
    *;

