package db

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/taha-ahmadi/simple-bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	argAccount1 := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: "USD",
	}
	argAccount2 := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: "USD",
	}
	account1, _ := testQueries.CreateAccount(context.Background(), argAccount1)
	account2, _ := testQueries.CreateAccount(context.Background(), argAccount2)

	count := 5
	amount, _ := decimal.NewFromString("10.0000")
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
		require.Equal(t, amount.StringFixed(4), transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check transaction
		fromTransaction := result.FromTransaction
		require.NotEmpty(t, fromTransaction)
		require.Equal(t, account1.ID, fromTransaction.AccountID)
		require.Equal(t, amount.Neg().StringFixed(4), fromTransaction.Amount)
		require.NotZero(t, fromTransaction.CreatedAt)

		_, err = store.GetTransactionById(context.Background(), fromTransaction.ID)
		require.NoError(t, err)

		toTransaction := result.ToTransaction
		require.NotEmpty(t, fromTransaction)
		require.Equal(t, account2.ID, toTransaction.AccountID)
		require.Equal(t, amount.StringFixed(4), toTransaction.Amount)
		require.NotZero(t, toTransaction.CreatedAt)

		_, err = store.GetTransactionById(context.Background(), toTransaction.ID)
		require.NoError(t, err)

		// check account's
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// check account's balance
		account1, _ := decimal.NewFromString(account1.Balance)
		account2, _ := decimal.NewFromString(account2.Balance)
		fromAccountBalance, _ := decimal.NewFromString(fromAccount.Balance)
		toAccountBalance, _ := decimal.NewFromString(toAccount.Balance)
		diff1 := account1.Sub(fromAccountBalance)
		diff2 := toAccountBalance.Sub(account2)
		require.Equal(t, diff1.String(), diff2.String())
		require.True(t, diff1.GreaterThan(decimal.NewFromInt(0)))
		r := diff1.Mod(amount)
		require.True(t, r.Equal(decimal.NewFromInt(0)))
	}

	// Check the final updated balances
	checkUpdatedBalances(t, account1, account2, int64(count), amount)

}
func checkUpdatedBalances(t *testing.T, account1, account2 Account, count int64, amount decimal.Decimal) {
	updatedAccount1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	decimalAccount1, _ := decimal.NewFromString(account1.Balance)
	decimalAccount2, _ := decimal.NewFromString(account2.Balance)
	n := decimal.NewFromInt(count)

	require.Equal(t, decimalAccount1.Sub(n.Mul(amount)).StringFixed(4), updatedAccount1.Balance)
	require.Equal(t, decimalAccount2.Add(n.Mul(amount)).StringFixed(4), updatedAccount2.Balance)
}
