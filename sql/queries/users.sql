-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: UpdateUserInfo :one
UPDATE users
SET email = $1, password = $2
WHERE id = $3
RETURNING id, created_at, updated_at, email;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
RETURNING *;

-- name: UserLogin :one
SELECT * FROM users WHERE email = $1;