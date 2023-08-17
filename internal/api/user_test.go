package api

import (
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/util"
	"testing"
)

func randomUser(t *testing.T) (user sql.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	f := faker.New()
	username := fmt.Sprintf("%s%s", util.RandomString(3), f.Internet().User())
	email := fmt.Sprintf("%s%s", util.RandomString(3), f.Internet().Email())

	user = sql.User{
		Username: username,
		Password: hashedPassword,
		Name:     f.Person().Name(),
		Email:    email,
	}
	return
}
