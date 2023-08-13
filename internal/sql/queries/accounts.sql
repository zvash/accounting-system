-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountById :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccountByIdForUpdate :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1 FOR NO KEY UPDATE;

-- name: GetAllAccounts :many
SELECT *
FROM accounts
ORDER BY created_at;

-- name: GetTwoAccountsInvolvedInTransfer :many
SELECT *
FROM accounts
WHERE id = @from_account
   OR id = @to_account
ORDER BY id
LIMIT 2 FOR NO KEY UPDATE;

-- name: GetAllAccountsPaginated :many
SELECT *
FROM accounts
ORDER BY created_at
OFFSET $1 LIMIT $2;

-- name: UpdateAccountById :one
UPDATE accounts
SET owner    = COALESCE(sqlc.narg(owner)::varchar, owner),
    balance  = COALESCE(sqlc.narg(balance)::bigint, balance),
    currency = COALESCE(sqlc.narg(currency)::varchar, currency)
WHERE id = sqlc.arg(id)::bigint
RETURNING *;

-- name: UpdateAccountBalanceById :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddToAccountBalanceById :one
UPDATE accounts
SET balance = balance + @amount
WHERE id = @id
RETURNING *;

-- name: DeleteAccountById :execrows
DELETE
FROM accounts
WHERE id = $1;