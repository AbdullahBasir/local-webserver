-- name: GetChirp :one
SELECT id FROM chirps WHERE id = $1;