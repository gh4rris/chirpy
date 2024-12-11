-- name: CreateToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    CURRENT_TIMESTAMP + INTERVAL '60 days',
    NULL
);

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
updated_at = NOW()
WHERE token = $1;