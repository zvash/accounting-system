-- name: CreateUser :one
INSERT INTO users (username, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUserName :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: UpdateUserByUsername :one
UPDATE users
SET email = COALESCE(sqlc.narg(email)::varchar, email),
    name  = COALESCE(sqlc.narg(name)::varchar, name)
WHERE username = sqlc.arg(username)
RETURNING *;