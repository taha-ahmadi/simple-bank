-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES ($1,$2,$3) RETURNING *;

-- name: GetTransferById :one
SELECT * FROM transfers 
WHERE id = $1 limit 1;

-- name: ListFromTransfers :many
SELECT * FROM transfers
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListToTransfers :many
SELECT * FROM transfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
