package api

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/brian926/UrlShorterGo/server/packages/config"
	"github.com/brian926/UrlShorterGo/server/packages/db"
)

func StartServer() *sql.DB {
	conn, err := db.ConnectDB()
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("DB connection error occurred")
	}
	defer conn.Close()

	runMigration := config.Config[config.RUN_MIGRATION]
	dbName := config.Config[config.POSTGRES_DB]

	if runMigration == "true" && conn != nil {
		if err := db.Migrate(conn, dbName); err != nil {
			log.WithField("reason", err.Error()).Fatal("db migration failed")
		}
	}

	return &conn
}
