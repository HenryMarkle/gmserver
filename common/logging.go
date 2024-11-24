package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Logger             *log.Logger
	DbConnectionString string
)

func init() {
	Logger = log.Default()

	_ = godotenv.Load()

	dbConnStr, dbConnFound := os.LookupEnv("DB_URL")

	if !dbConnFound {
		Logger.Fatal("Could not find database connection string")
	}

	DbConnectionString = dbConnStr
}
