package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Setup() error {
	dsn := "root:@tcp(10.39.50.237:3306)/nutrishe"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error verifying connection to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
	return nil
}

func GetDB() *sql.DB {
	return db
}
