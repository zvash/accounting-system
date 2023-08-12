-- name: CreateEntry :one
INSERT INTO entries (account_id, amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetEntryById :one
SELECT *
FROM entries
WHERE id = $1
LIMIT 1;

-- name: GetAllEntries :many
SELECT *
FROM entries
ORDER BY created_at;

-- name: GetAllEntriesPaginated :many
SELECT *
FROM entries
ORDER BY created_at
OFFSET $1 LIMIT $2;

-- name: UpdateEntryById :one
UPDATE entries
SET account_id = coalesce(sqlc.narg('account_id'), account_id),
    amount     = coalesce(sqlc.narg('amount'), amount)
WHERE id = $1
RETURNING *;

-- name: DeleteEntryById :execrows
DELETE
FROM entries
WHERE id = $1;