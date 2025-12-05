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
	"home/osarukun/repos/tower-investing/backend/api"

	db "home/osarukun/repos/tower-investing/backend/db/sqlc"
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
		log.Fatal("cannot connect to DB: ", err)
	}

	conn.Exec("PRAGMA journal_mode=WAL;")
	conn.Exec("PRAGMA foreign_keys=ON;")

	err = migration.RunMigrations(conn, "./db/migration/")
	if err != nil {
		log.Fatal("migration error: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
