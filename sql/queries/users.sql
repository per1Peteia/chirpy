-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
    VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2)
RETURNING
    *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: UpdateUserCreds :one
UPDATE
    users
SET
    hashed_password = $1,
    email = $2,
    updated_at = now()
WHERE
    id = $3
RETURNING
    *;

-- name: UpgradeUserMembershipByID :exec
UPDATE
    users
SET
    is_chirpy_red = TRUE
WHERE
    id = $1;

