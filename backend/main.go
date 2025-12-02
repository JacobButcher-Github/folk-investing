package main

import (
	//stl
	"database/sql"
	"log"

	//go package

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/api"
	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
)

const (
	dbDriver      = "sqlite"
	dbSource      = ""
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connec to DB: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
