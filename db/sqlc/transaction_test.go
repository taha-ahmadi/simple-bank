package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/taha-ahmadi/simple-bank/util"
)

func createRandomTransaction(t *testing.T, account Account) Transaction {
	arg := CreateTransactionParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.AccountID, transaction.AccountID)
	require.Equal(t, arg.Amount, transaction.Amount)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	arg := CreateTransactionParams{
		AccountID: 1,
		Amount:    util.RandomBalance(),
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.Amount, transaction.Amount)
	require.Equal(t, arg.AccountID, transaction.AccountID)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)
}

func TestGetTransaction(t *testing.T) {
	arg := CreateTransactionParams{
		AccountID: 1,
		Amount:    util.RandomBalance(),
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	getTransaction, err := testQueries.GetTransactionById(context.Background(), transaction.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, transaction.AccountID, getTransaction.AccountID)
	require.Equal(t, transaction.Amount, getTransaction.Amount)
	require.WithinDuration(t, transaction.CreatedAt, getTransaction.CreatedAt, time.Second)

}

func TestListTransactions(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		arg := CreateTransactionParams{
			AccountID: account.ID,
			Amount:    util.RandomBalance(),
		}
	
		_, err := testQueries.CreateTransaction(context.Background(), arg)
		if err != nil {
			fmt.Println(err)
		}
	}

	arg := ListAccountTransactionsParams{
		AccountID: account.ID,
		Limit:  5,
		Offset: 5,
	}

	transactions, err := testQueries.ListAccountTransactions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
	}
}
