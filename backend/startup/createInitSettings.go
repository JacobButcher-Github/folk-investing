package startup

import (
	"context"
	"database/sql"
	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
	"github.com/JacobButcher-Github/folk-investing/backend/util"
	"time"
)

func CreateSettings(config util.Config, store db.Store) error {
	ctx := context.Background()

	// Check if settings exists
	_, err := store.GetSiteSettings(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	// Create initial settings
	_, err = store.CreateSiteSettings(ctx, db.CreateSiteSettingsParams{
		NumberOfEventsVisible: 10,
		ValueSymbol:           "$",
		EventLabel:            "Instance",
		LockoutTimeStart:      time.Now(),
	})

	return err
}
