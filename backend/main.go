package main

import (
	//stl

	"database/sql"
	migration "home/osarukun/repos/tower-investing/backend/db"
	"home/osarukun/repos/tower-investing/backend/util"
	"log"

	//go package
	_ "modernc.org/sqlite"

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/api"
	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
)

const (
	dbDriver = "sqlite"
	dbSource = "./db/app.db"
)

func main() {
	config, err := util.ReadConfig()
	if err != nil {
		log.Fatal("could not read config: ", err)
	}

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connec to DB: ", err)
	}

	conn.Exec("PRAGMA journal_mode=WAL;")

	err = migration.RunMigrations(conn, "./db/migration/")

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
