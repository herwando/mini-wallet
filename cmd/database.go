package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseClientOption struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func getDatabaseDefaultOption() DatabaseClientOption {
	defaultOpt := DatabaseClientOption{
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		Database: os.Getenv("PSQL_DATABASE"),
	}

	return defaultOpt
}

func getDBConnection() *sql.DB {
	options := getDatabaseDefaultOption()
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", options.Host, options.Port, options.User, options.Password, options.Database)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("mini-wallet database error:", err)
		panic(err)
	}

	return db
}
