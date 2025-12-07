package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
    `)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", name).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			continue // already applied
		}

		// read migration file
		content, err := os.ReadFile(filepath.Join(migrationsPath, name))
		if err != nil {
			return err
		}

		// apply SQL
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("migration %s failed: %w", name, err)
		}

		// mark migration as applied
		_, err = db.Exec("INSERT INTO migrations(name) VALUES (?)", name)
		if err != nil {
			return err
		}

		fmt.Println("Applied migration:", name)
	}

	return nil
}
