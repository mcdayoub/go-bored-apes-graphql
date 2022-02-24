-- name: ListTransfersByTransaction :many
SELECT * FROM transfers
WHERE transaction = $1;

-- name: ListTransfersBySender :many
SELECT * FROM transfers
WHERE sender = $1;

-- name: ListTransfersByReceiver :many
SELECT * FROM transfers
WHERE receiver = $1;

-- name: ListTransfersByTokenID :many
SELECT * FROM transfers
WHERE token_id = $1;

-- name: ListUnreadTransfers :many
SELECT * FROM transfers
Where read = FALSE;

-- name: CreateTransfer :one
INSERT INTO transfers (transaction, sender, receiver, token_id, read)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ReadTransfer :one
UPDATE transfers
SET read = TRUE
WHERE transaction = $1
RETURNING *;