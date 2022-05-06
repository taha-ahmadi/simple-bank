-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES ($1,$2,$3) RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts 
WHERE id = $1 limit 1;

-- name: GetAccountByIdForUpdate :one
SELECT * FROM accounts
WHERE id = $1 limit 1 FOR NO KEY UPDATE ;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM accounts
where id = $1
RETURNING id;