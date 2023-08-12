package sql

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/zvash/accounting-system/internal/util"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, args.Owner)
	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Currency, args.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func getRandomlyCreatedAccountById(t *testing.T, id int64) Account {
	account, err := testQueries.GetAccountById(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	return account
}

func TestCreateAccount(t *testing.T) {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, args.Owner)
	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Currency, args.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccountById(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestUpdateAccountById(t *testing.T) {
	account1 := createRandomAccount(t)
	owner := util.RandomOwner()
	args := UpdateAccountByIdParams{
		ID:       account1.ID,
		Owner:    sql.NullString{String: owner, Valid: true},
		Balance:  sql.NullInt64{},
		Currency: sql.NullString{},
	}
	affected, err := testQueries.UpdateAccountById(context.Background(), args)
	require.NoError(t, err)
	require.NotZero(t, affected)

	account2 := getRandomlyCreatedAccountById(t, account1.ID)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, owner, account2.Owner)
	require.NotEqual(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestDeleteAccountById(t *testing.T) {
	account := createRandomAccount(t)
	affected, err := testQueries.DeleteAccountById(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotZero(t, affected)

	account, err = testQueries.GetAccountById(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestGetAllAccountsPaginated(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	args := GetAllAccountsPaginatedParams{
		Offset: 5,
		Limit:  5,
	}
	accounts, err := testQueries.GetAllAccountsPaginated(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
