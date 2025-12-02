package db

import (
	//stl
	"context"
	"testing"
	"time"

	//go package
	"github.com/stretchr/testify/require"

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/util"
)

func createRandomSiteSettings(t *testing.T) SiteSetting {
	arg := CreateSiteSettingsParams{
		NumberOfEventsVisible: util.RandomInt(1, 50),
		ValueSymbol:           util.RandomString(1),
		EventLabel:            util.RandomString(10),
		LockoutTimeStart:      time.Now(),
	}

	settings, err := testQueries.CreateSiteSettings(context.Background(), arg)

	require.Empty(t, err)
	require.NotEmpty(t, settings)

}
