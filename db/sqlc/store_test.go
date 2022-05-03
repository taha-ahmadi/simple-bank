package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	count := 5
	amount := "10.0000"
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < count; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check Result
	for i := 0; i < count; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check transactoin
		fromTransaction := result.FromTransaction
		require.NotEmpty(t, fromTransaction)
		require.Equal(t, account1.ID, fromTransaction.AccountID)
		require.Equal(t, "-" + amount, fromTransaction.Amount)
		require.NotZero(t, fromTransaction.CreatedAt)

		_, err = store.GetTransactionById(context.Background(), fromTransaction.ID)
		require.NoError(t, err)

		toTransaction := result.ToTransaction
		require.NotEmpty(t, fromTransaction)
		require.Equal(t, account2.ID, toTransaction.AccountID)
		require.Equal(t, amount, toTransaction.Amount)
		require.NotZero(t, toTransaction.CreatedAt)

		_, err = store.GetTransactionById(context.Background(), toTransaction.ID)
		require.NoError(t, err)
	}
}
