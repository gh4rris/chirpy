-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens
ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL AND expires_at > NOW();

-- name: UpdateUserEmailPassword :one
UPDATE users
SET email = $2, hashed_password = $3,
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;