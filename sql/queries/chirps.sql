-- name: CreateChirps :one
INSERT INTO chirps(
    id, 
    created_at, 
    updated_at, 
    body, 
    user_id
) VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * 
FROM chirps
ORDER BY created_at ASC;

-- name: GetSingleChirp :one
SELECT * 
FROM chirps
WHERE chirps.id = $1
LIMIT 1;

-- name: GetUserChirps :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE chirps.id = $1;
