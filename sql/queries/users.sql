-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE
FROM users;

-- name: GetUserId :one
SELECT id 
FROM users
WHERE email = $1;

-- name: GetUserHash :one
SELECT hashed_password
FROM users
WHERE email = $1;

-- name: GetUserFromEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT *
FROM users
WHERE id = (
SELECT user_id
FROM refresh_tokens
WHERE token = $1);

-- name: UpdateUserData :one
UPDATE users
SET hashed_password = $1,
email = $2
WHERE id = $3
RETURNING *;

-- name: SetIsRed :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;