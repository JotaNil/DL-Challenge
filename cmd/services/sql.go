package services

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"os"
)

const (
	// Ipv4ProxyDB ipv4-proxy-dreamlab DB
	Ipv4ProxyDB = "ipv4ProxyDB"
	// MockDB for unit test usage
	MockDB = "mock"

	ipv4ProxyHost   = "localhost"
	ipv4ProxyPort   = 5432
	ipv4ProxyUser   = "postgres"
	ipv4ProxyDbname = "ipv4-proxy-dreamlab"
)

var ipv4ProxyPassword = os.Getenv("DL_CHALLENGE_DBPASS")

func ConnectToSQLDB(dbName string) (*sql.DB, sqlmock.Sqlmock) {
	switch dbName {
	case Ipv4ProxyDB:
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			ipv4ProxyHost, ipv4ProxyPort, ipv4ProxyUser, ipv4ProxyPassword, ipv4ProxyDbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		return db, nil
	case "mock":
		db, mock, _ := sqlmock.New()
		return db, mock
	default:
		panic(fmt.Sprintf("db %s is not declared in connections", dbName))
	}
}
