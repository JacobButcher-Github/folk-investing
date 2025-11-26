package db

import (
	//stl
	"context"
	"database/sql"
	"testing"

	//go package
	"github.com/stretchr/testify/require"

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/util"
)

func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		UserLogin:      util.RandomString(6),
		HashedPassword: util.RandomString(12),
		DisplayName:    sql.NullString{String: util.RandomString(5), Valid: true},
		Dollars:        sql.NullInt64{Int64: util.RandomInt(3, 10), Valid: true},
		Cents:          sql.NullInt64{Int64: util.RandomInt(0, 10), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserLogin, user.UserLogin)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.DisplayName, user.DisplayName)
	require.Equal(t, arg.Dollars, user.Dollars)
	require.Equal(t, arg.Cents, user.Cents)

	require.NotZero(t, user.ID)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.UserLogin) {

	}
}
