// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: transaction.sql

package db

import (
	"context"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    account_id,
    amount
) VALUES ($1,$2) RETURNING id, account_id, amount, created_at
`

type CreateTransactionParams struct {
	AccountID int64  `json:"account_id"`
	Amount    string `json:"amount"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.queryRow(ctx, q.createTransactionStmt, createTransaction, arg.AccountID, arg.Amount)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransactionById = `-- name: GetTransactionById :one
SELECT id, account_id, amount, created_at FROM transactions 
WHERE id = $1 limit 1
`

func (q *Queries) GetTransactionById(ctx context.Context, id int64) (Transaction, error) {
	row := q.queryRow(ctx, q.getTransactionByIdStmt, getTransactionById, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listAccountTransactions = `-- name: ListAccountTransactions :many
SELECT id, account_id, amount, created_at FROM transactions
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListAccountTransactionsParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListAccountTransactions(ctx context.Context, arg ListAccountTransactionsParams) ([]Transaction, error) {
	rows, err := q.query(ctx, q.listAccountTransactionsStmt, listAccountTransactions, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
