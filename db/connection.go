package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"os"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	certPool := x509.NewCertPool()
	pemBytes, pemErr := os.ReadFile(`./ca.pem`)
	if pemErr != nil {
		common.Logger.Fatalf("Failed to load certificate: %v\n", pemErr)
		return
	}

	certPool.AppendCertsFromPEM(pemBytes)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            certPool,
	}

	tlsErr := mysql.RegisterTLSConfig("required", tlsConfig)

	mysql.RegisterTLSConfig("custom", tlsConfig)

	if tlsErr != nil {
		common.Logger.Printf("Failed to configure database TLS: %v\n", tlsErr)
		return
	}

	db, dbErr := sql.Open("mysql", common.DbConnectionString)

	if dbErr != nil {
		common.Logger.Fatalf("Failed to connect to database: %v\n", dbErr)
	}

	if err := db.Ping(); err != nil {
		common.Logger.Fatalf("Failed to ping database: %v", err)
	}

	DB = db
}
