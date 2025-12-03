package db

import (
	//stl
	"context"
	"database/sql"
	"testing"
	"time"

	//go package
	"github.com/stretchr/testify/require"

	//local
	"home/osarukun/repos/tower-investing/backend/util"
)

func createRandomSiteSettings(t *testing.T) SiteSetting {
	arg := CreateSiteSettingsParams{
		NumberOfEventsVisible: util.RandomInt(1, 50),
		ValueSymbol:           util.RandomString(1),
		EventLabel:            util.RandomString(10),
		LockoutTimeStart:      time.Now(),
	}

	settings, err := testQueries.CreateSiteSettings(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, settings)

	require.Equal(t, settings.NumberOfEventsVisible, arg.NumberOfEventsVisible)
	require.Equal(t, settings.ValueSymbol, arg.ValueSymbol)
	require.Equal(t, settings.EventLabel, arg.EventLabel)
	require.Equal(t, settings.LockoutTimeStart, arg.LockoutTimeStart)
	return settings
}

func TestCreateSiteSettings(t *testing.T) {
	//only call once at the beginning of tests, then use update for everything else.
	createRandomSiteSettings(t)
}

func TestGetSiteSettings(t *testing.T) {
	settings1, err := testQueries.GetSiteSettings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, settings1)

	settings2, err := testQueries.GetSiteSettings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, settings2)

	require.Equal(t, settings1.NumberOfEventsVisible, settings2.NumberOfEventsVisible)
	require.Equal(t, settings1.ValueSymbol, settings2.ValueSymbol)
	require.Equal(t, settings1.EventLabel, settings2.EventLabel)
	require.Equal(t, settings1.LockoutTimeStart, settings2.LockoutTimeStart)
}

func TestUpdateSiteSettingsNumEvents(t *testing.T) {
	settings1, err := testQueries.GetSiteSettings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, settings1)

	updatedSettingsArgs := UpdateSettingsParams{
		NumberOfEventsVisible: sql.NullInt64{Int64: util.RandomInt(1, 50), Valid: true},
		ValueSymbol:           sql.NullString{String: "", Valid: false},
		EventLabel:            sql.NullString{String: "", Valid: false},
		LockoutTimeStart:      sql.NullTime{Time: time.Now(), Valid: false},
	}
	updatedSettings, err := testQueries.UpdateSettings(context.Background(), updatedSettingsArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSettings)

	require.Equal(t, updatedSettingsArgs.NumberOfEventsVisible, updatedSettings.NumberOfEventsVisible)
	require.NotEqual(t, settings1.NumberOfEventsVisible, updatedSettings.NumberOfEventsVisible)

	require.Equal(t, settings1.ValueSymbol, updatedSettings.ValueSymbol)
	require.Equal(t, settings1.EventLabel, updatedSettings.EventLabel)
	require.Equal(t, settings1.LockoutTimeStart, updatedSettings.LockoutTimeStart)
}
func TestUpdateSiteSettingsValueSymbol(t *testing.T) {
	settings1, err := testQueries.GetSiteSettings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, settings1)

	updatedSettingsArgs := UpdateSettingsParams{
		NumberOfEventsVisible: sql.NullInt64{Int64: 0, Valid: false},
		ValueSymbol:           sql.NullString{String: util.RandomString(1), Valid: true},
		EventLabel:            sql.NullString{String: "", Valid: false},
		LockoutTimeStart:      sql.NullTime{Time: time.Now(), Valid: false},
	}
	updatedSettings, err := testQueries.UpdateSettings(context.Background(), updatedSettingsArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSettings)

	require.Equal(t, updatedSettingsArgs.ValueSymbol, updatedSettings.ValueSymbol)
	require.NotEqual(t, settings1.ValueSymbol, updatedSettings.ValueSymbol)

	require.Equal(t, settings1.NumberOfEventsVisible, updatedSettings.NumberOfEventsVisible)
	require.Equal(t, settings1.EventLabel, updatedSettings.EventLabel)
	require.Equal(t, settings1.LockoutTimeStart, updatedSettings.LockoutTimeStart)
}
func TestUpdateSiteSettingsEventLabel(t *testing.T) {
	settings1, err := testQueries.GetSiteSettings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, settings1)

	updatedSettingsArgs := UpdateSettingsParams{
		NumberOfEventsVisible: sql.NullInt64{Int64: 0, Valid: false},
		ValueSymbol:           sql.NullString{String: "", Valid: false},
		EventLabel:            sql.NullString{String: util.RandomString(10), Valid: true},
		LockoutTimeStart:      sql.NullTime{Time: time.Now(), Valid: false},
	}
	updatedSettings, err := testQueries.UpdateSettings(context.Background(), updatedSettingsArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSettings)

	require.Equal(t, updatedSettingsArgs.EventLabel, updatedSettings.EventLabel)
	require.NotEqual(t, settings1.EventLabel, updatedSettings.EventLabel)

	require.Equal(t, settings1.NumberOfEventsVisible, updatedSettings.NumberOfEventsVisible)
	require.Equal(t, settings1.ValueSymbol, updatedSettings.ValueSymbol)
	require.Equal(t, settings1.LockoutTimeStart, updatedSettings.LockoutTimeStart)
}
func TestUpdateSiteSettingsLockout(t *testing.T) {
	settings1, err := testQueries.GetSiteSettings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, settings1)

	updatedSettingsArgs := UpdateSettingsParams{
		NumberOfEventsVisible: sql.NullInt64{Int64: 0, Valid: false},
		ValueSymbol:           sql.NullString{String: "", Valid: false},
		EventLabel:            sql.NullString{String: "", Valid: false},
		LockoutTimeStart:      sql.NullTime{Time: time.Now(), Valid: true},
	}
	updatedSettings, err := testQueries.UpdateSettings(context.Background(), updatedSettingsArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSettings)

	require.Equal(t, updatedSettingsArgs.LockoutTimeStart, updatedSettings.LockoutTimeStart)
	require.NotEqual(t, settings1.LockoutTimeStart, updatedSettings.LockoutTimeStart)

	require.Equal(t, settings1.NumberOfEventsVisible, updatedSettings.NumberOfEventsVisible)
	require.Equal(t, settings1.ValueSymbol, updatedSettings.ValueSymbol)
	require.Equal(t, settings1.EventLabel, updatedSettings.EventLabel)
}

func TestUpdateSiteSettingsAllFields(t *testing.T) {
	settings1, err := testQueries.GetSiteSettings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, settings1)

	updatedSettingsArgs := UpdateSettingsParams{
		NumberOfEventsVisible: sql.NullInt64{Int64: util.RandomInt(1, 50), Valid: true},
		ValueSymbol:           sql.NullString{String: util.RandomString(1), Valid: true},
		EventLabel:            sql.NullString{String: util.RandomString(10), Valid: true},
		LockoutTimeStart:      sql.NullTime{Time: time.Now(), Valid: true},
	}
	updatedSettings, err := testQueries.UpdateSettings(context.Background(), updatedSettingsArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSettings)

	require.Equal(t, updatedSettingsArgs.NumberOfEventsVisible, updatedSettings.NumberOfEventsVisible)
	require.NotEqual(t, settings1.NumberOfEventsVisible, updatedSettings.NumberOfEventsVisible)
	require.Equal(t, updatedSettingsArgs.ValueSymbol, updatedSettings.ValueSymbol)
	require.NotEqual(t, settings1.ValueSymbol, updatedSettings.ValueSymbol)
	require.Equal(t, updatedSettingsArgs.EventLabel, updatedSettings.EventLabel)
	require.NotEqual(t, settings1.EventLabel, updatedSettings.EventLabel)
	require.Equal(t, updatedSettingsArgs.LockoutTimeStart, updatedSettings.LockoutTimeStart)
	require.NotEqual(t, settings1.LockoutTimeStart, updatedSettings.LockoutTimeStart)
}
