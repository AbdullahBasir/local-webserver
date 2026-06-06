-- name: DeleteChirp :exec
Delete FROM chirps WHERE id = $1;