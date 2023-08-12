-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountById :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAllAccounts :many
SELECT *
FROM accounts
ORDER BY created_at;

-- name: GetAllAccountsPaginated :many
SELECT *
FROM accounts
ORDER BY created_at
OFFSET $1 LIMIT $2;

-- name: UpdateAccountById :one
UPDATE accounts
SET owner    = coalesce(sqlc.narg('owner'), owner),
    balance  = coalesce(sqlc.narg('balance'), balance),
    currency = coalesce(sqlc.narg('currency'), currency)
WHERE id = $1
RETURNING *;

-- name: DeleteAccountById :execrows
DELETE
FROM accounts
WHERE id = $1;