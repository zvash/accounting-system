package sql

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTransaction(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//run n concurrent transfer transactions to make sure it can
	//handle concurrent calls correctly
	n := 5
	var amount int64 = 10

	errorsChannel := make(chan error)
	resultsChannel := make(chan TransferTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errorsChannel <- err
			resultsChannel <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errorsChannel
		require.NoError(t, err)

		result := <-resultsChannel
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO: check accounts' balance
	}
}
