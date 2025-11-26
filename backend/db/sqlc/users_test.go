package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		UserLogin:      "Chatter",
		HashedPassword: "abababa",
		DisplayName:    sql.NullString{String: "Chatt", Valid: true},
		Dollars:        sql.NullInt64{Int64: 15, Valid: true},
		Cents:          sql.NullInt64{Int64: 15, Valid: true},
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
}
