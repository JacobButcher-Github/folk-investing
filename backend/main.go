package main

import (
	//stl

	"database/sql"
	migration "github.com/JacobButcher-Github/folk-investing/backend/db"
	"github.com/JacobButcher-Github/folk-investing/backend/startup"
	"github.com/JacobButcher-Github/folk-investing/backend/util"
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
		log.Fatal("cannot connect to DB: ", err)
	}

	conn.Exec("PRAGMA busy_timeout = 5000;")
	conn.Exec("PRAGMA journal_mode=WAL;")
	conn.Exec("PRAGMA foreign_keys=ON;")

	err = migration.RunMigrations(conn, "./db/migration/")
	if err != nil {
		log.Fatal("migration error: ", err)
	}

	store := db.NewStore(conn)

	//Initial admin account and settings
	startup.CreateInitialAdmin(config, *store)
	startup.CreateSettings(config, *store)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
