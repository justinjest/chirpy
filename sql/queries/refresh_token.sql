-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    now() + INTERVAL '60 days', 
    NULL
)
RETURNING *;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET updated_at = now(),
revoked_at = now()
WHERE user_id = $1
RETURNING *;

-- name: SelectNewestToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1
AND revoked_at IS NULL
AND expires_at > NOW();