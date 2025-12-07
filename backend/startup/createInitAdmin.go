package startup

import (
	"context"
	"database/sql"
	db "home/osarukun/repos/tower-investing/backend/db/sqlc"
	"home/osarukun/repos/tower-investing/backend/util"
)

func CreateInitialAdmin(config util.Config, store db.Store) error {
	ctx := context.Background()

	// Check if admin exists
	_, err := store.GetUserFromName(ctx, config.AdminUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	// Create initial admin
	hashedPassword, err := util.HashPassword(config.AdminPassword)
	if err != nil {
		return err
	}
	_, err = store.CreateAdmin(ctx, db.CreateAdminParams{
		UserLogin:      config.AdminUsername,
		HashedPassword: hashedPassword,
		Role:           util.AdminRole,
		Dollars:        100,
		Cents:          0,
	})

	return err
}
