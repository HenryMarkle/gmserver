package db

import (
	"database/sql"

	"github.com/HenryMarkle/gmserver/common"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	db, dbErr := sql.Open("mysql", common.DbConnectionString)

	if dbErr != nil {
		common.Logger.Fatalf("Failed to connect to database: %v\n", dbErr)
	}

	DB = db
}
