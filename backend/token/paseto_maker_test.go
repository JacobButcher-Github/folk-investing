package token

import (
	//stl

	"errors"
	"testing"
	"time"

	//go package

	"github.com/stretchr/testify/require"

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	userLogin := util.RandomString(15)
	userID := util.RandomInt(1, 100)
	role := util.UserRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userID, userLogin, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userLogin, payload.UserLogin)
	require.Equal(t, userID, payload.UserID)
	require.Equal(t, role, payload.Role)
	require.Equal(t, issuedAt, payload.IssuedAt, time.Second)
	require.Equal(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomInt(1, 100), util.RandomString(15), util.UserRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errors.New("token is expired").Error())
	require.Nil(t, payload)
}
