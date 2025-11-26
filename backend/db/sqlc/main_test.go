package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

const (
	dbDriver = "sqlite"
	dbSource = ""
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
