-- name: UpdateUserInfo :one
UPDATE users
SET email = $1, password = $2
WHERE id = $3
RETURNING id, created_at, updated_at, email;