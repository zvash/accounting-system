package sql

import (
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
	"github.com/zvash/accounting-system/internal/util"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	f := faker.New()
	passwordString := util.RandomString(8)
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(passwordString), 10)

	username := fmt.Sprintf("%s%s", util.RandomString(3), f.Internet().User())
	email := fmt.Sprintf("%s%s", util.RandomString(3), f.Internet().Email())
	require.NoError(t, err)
	args := CreateUserParams{
		Username: username,
		Name:     f.Person().Name(),
		Email:    email,
		Password: string(passwordBytes),
	}
	user, err := testStore.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.Name, user.Name)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.Password, user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordString))
	require.NoError(t, err)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByUserName(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUserByUserName(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
