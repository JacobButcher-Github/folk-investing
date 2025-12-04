package api

import (
	"database/sql"
	db "home/osarukun/repos/tower-investing/backend/db/sqlc"
	"home/osarukun/repos/tower-investing/backend/util"
	"log"
	"os"
	"testing"
	"time"

	migration "home/osarukun/repos/tower-investing/backend/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var testStore *db.Store

func NewTestServer(t *testing.T) *Server {

	config := util.Config{
		AdminUsername:        "tester",
		AdminPassword:        "test",
		ServerAddress:        "127.0.0.1:8080",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 24 * time.Hour,
	}

	server, err := NewServer(config, testStore)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	conn, err := sql.Open("sqlite", "./db/test-app.db")
	if err != nil {
		log.Fatal("cannot connect to DB:", err)
	}

	// WAL mode (same as production)
	conn.Exec("PRAGMA journal_mode=WAL;")

	// Run migrations
	err = migration.RunMigrations(conn, "./db/migration/")
	if err != nil {
		log.Fatal("cannot run migrations:", err)
	}

	testStore = db.NewStore(conn)

	os.Exit(m.Run())
}
