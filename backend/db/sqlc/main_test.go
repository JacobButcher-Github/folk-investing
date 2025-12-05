package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	migration "home/osarukun/repos/tower-investing/backend/db"

	_ "modernc.org/sqlite"
)

const (
	dbDriver = "sqlite"
	dbSource = "../test-app.db"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	testDB.Exec("PRAGMA journal_mode=WAL;")
	testDB.Exec("PRAGMA foreign_keys = ON;")

	err = migration.RunMigrations(testDB, "../migration/")
	if err != nil {
		log.Fatal("migration error: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
