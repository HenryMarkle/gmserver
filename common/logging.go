package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Logger             *log.Logger
	DbConnectionString string

	StoragePath string
)

func init() {
	Logger = log.Default()

	_ = godotenv.Load()

	dbConnStr, dbConnFound := os.LookupEnv("DB_URL")

	if !dbConnFound {
		Logger.Fatal("Could not find database connection string")
		os.Exit(-1)
	}

	DbConnectionString = dbConnStr

	storagePath, pathFound := os.LookupEnv("STORAGE_PATH")
	if !pathFound {
		// Logger.Fatal("STORAGE_PATH not set")
		// os.Exit(-1)
		storagePath = "."
	}

	StoragePath = storagePath
}
