-- name: CreateUser :one
INSERT INTO users (username, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUserName :one
SELECT * FROM users WHERE username = $1 LIMIT 1;