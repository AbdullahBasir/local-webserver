-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens 
WHERE token = $1 AND NOW() < expires_at AND revoked_at IS NULL;