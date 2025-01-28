-- name: CreateUser :one
INSERT INTO users (
    id, 
    created_at, 
    updated_at, 
    email, 
    hashed_password
) VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE users.email = $1
LIMIT 1;

-- name: UpdateUsersCredentials :one
UPDATE users
SET hashed_password = $1, email = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email;