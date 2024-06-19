package models

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/go-sql-driver/mysql"
// )

// var db *sql.DB

// func Setup() error {
// 	dsn := "root:localhost@tcp(172.0.0.1:3306)/empowher"

// 	var err error
// 	db, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		return fmt.Errorf("error opening database: %v", err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		return fmt.Errorf("error verifying connection to the database: %v", err)
// 	}

// 	fmt.Println("Successfully connected to the database")
// 	return nil
// }

// // GetDB returns a reference to the database
// func GetDB() *sql.DB {
// 	return db
// }

import (
	// "nutrishe/models"

	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	// "gorm.io/gorm"
)

var db *sql.DB

func Setup() error {
	dsn := "root:@tcp(192.168.1.7:3306)/nutrishe"

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

// func ConnectDatabase() {
// 	database, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go-api-nutrishe"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	database.AutoMigrate(
// 		&models.DietPlan{})

// 	DB = database
// }
