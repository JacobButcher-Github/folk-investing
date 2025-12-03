package db

import (
	//stl
	"context"
	"database/sql"
	"testing"

	//go package
	"github.com/stretchr/testify/require"

	//local
	"home/osarukun/repos/tower-investing/backend/util"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		UserLogin:      util.RandomString(6),
		HashedPassword: hashedPassword,
		Dollars:        util.RandomInt(0, 100000),
		Cents:          util.RandomInt(0, 99),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserLogin, user.UserLogin)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
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
	user2, err := testQueries.GetUserFromName(context.Background(), user1.UserLogin)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Dollars, user2.Dollars)
	require.Equal(t, user1.Cents, user2.Cents)
}

func TestUpdateUserPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	var newHashedPassword string
	var err error

	for {
		newPassword := util.RandomString(6)
		newHashedPassword, err = util.HashPassword(newPassword)
		require.NoError(t, err)
		if oldUser.HashedPassword != newHashedPassword {
			break
		}
	}

	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		sql.NullString{String: newHashedPassword, Valid: true},
		sql.NullInt64{Int64: 0, Valid: false},
		sql.NullInt64{Int64: 0, Valid: false},
		oldUser.UserLogin,
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.UserLogin, updatedUser.UserLogin)
	require.Equal(t, oldUser.Dollars, updatedUser.Dollars)
	require.Equal(t, oldUser.Cents, updatedUser.Cents)
}

func TestUpdateUserDollar(t *testing.T) {
	oldUser := createRandomUser(t)
	newDollars := util.RandomInt(0, 100000)

	for {
		if oldUser.Dollars != newDollars {
			break
		}
		newDollars = util.RandomInt(0, 100000)
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		sql.NullString{String: "", Valid: false},
		sql.NullInt64{Int64: newDollars, Valid: true},
		sql.NullInt64{Int64: 0, Valid: false},
		oldUser.UserLogin,
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Dollars, updatedUser.Dollars)
	require.Equal(t, newDollars, updatedUser.Dollars)
	require.Equal(t, oldUser.UserLogin, updatedUser.UserLogin)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Cents, updatedUser.Cents)
}

func TestUpdateUserCent(t *testing.T) {
	oldUser := createRandomUser(t)
	newCents := util.RandomInt(0, 99)

	for {
		if oldUser.Cents != newCents {
			break
		}
		newCents = util.RandomInt(0, 99)
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		sql.NullString{String: "", Valid: false},
		sql.NullInt64{Int64: 0, Valid: false},
		sql.NullInt64{Int64: newCents, Valid: true},
		oldUser.UserLogin,
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Cents, updatedUser.Cents)
	require.Equal(t, newCents, updatedUser.Cents)
	require.Equal(t, oldUser.UserLogin, updatedUser.UserLogin)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Dollars, updatedUser.Dollars)
}

func TestUpdateUserMoney(t *testing.T) {
	oldUser := createRandomUser(t)
	newDollars := util.RandomInt(0, 100000)
	newCents := util.RandomInt(0, 99)

	for {
		if oldUser.Cents != newCents && oldUser.Dollars != newDollars {
			break
		}
		if oldUser.Cents == newCents {
			newCents = util.RandomInt(0, 99)
		} else {
			newDollars = util.RandomInt(0, 100000)
		}
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		sql.NullString{String: "", Valid: false},
		sql.NullInt64{Int64: newDollars, Valid: true},
		sql.NullInt64{Int64: newCents, Valid: true},
		oldUser.UserLogin,
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Dollars, updatedUser.Dollars)
	require.NotEqual(t, oldUser.Cents, updatedUser.Cents)
	require.Equal(t, newDollars, updatedUser.Dollars)
	require.Equal(t, newCents, updatedUser.Cents)
	require.Equal(t, oldUser.UserLogin, updatedUser.UserLogin)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newDollars := util.RandomInt(0, 100000)
	newCents := util.RandomInt(0, 99)

	for {
		if oldUser.Cents != newCents &&
			oldUser.Dollars != newDollars &&
			oldUser.HashedPassword != newHashedPassword {
			break
		}
		if oldUser.Cents == newCents {
			newCents = util.RandomInt(0, 99)
		} else if oldUser.Dollars == newDollars {
			newDollars = util.RandomInt(0, 100000)
		} else {
			newPassword = util.RandomString(6)
			newHashedPassword, err = util.HashPassword(newPassword)
		}
	}

	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		sql.NullString{String: newHashedPassword, Valid: true},
		sql.NullInt64{Int64: newDollars, Valid: true},
		sql.NullInt64{Int64: newCents, Valid: true},
		oldUser.UserLogin,
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.Dollars, updatedUser.Dollars)
	require.NotEqual(t, oldUser.Cents, updatedUser.Cents)
	require.Equal(t, oldUser.UserLogin, updatedUser.UserLogin)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newDollars, updatedUser.Dollars)
	require.Equal(t, newCents, updatedUser.Cents)
}
