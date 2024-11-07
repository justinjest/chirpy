-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE
FROM users;

-- name: GetUserId :one
SELECT id 
FROM users
WHERE email == $1;