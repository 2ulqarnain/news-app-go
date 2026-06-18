package database

import (
	"GoNewsScrapper/internals/config"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	if config.DbFilePath == "" {
		return nil, fmt.Errorf("no DB File Path Found")
	}

	db, err := sql.Open("sqlite3", config.DbFilePath)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
