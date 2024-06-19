package models

import (
	"database/sql"
	"log"
	"nutrishe/entity"
)

func CreateNewMeal(db *sql.DB, data entity.Food) error {
	_, err := db.Exec("INSERT INTO food (FoodID, Name, Serving, Calories, Type) VALUE (?, ?, ?, ?, 'food')", data.FoodID, data.Name, data.Serving, data.Calories)

	if err != nil {
		return err
	}

	return nil
}

func ViewMeal(db *sql.DB) entity.Food {
	row := db.QueryRow("SELECT FoodID, Name, Serving, Calories FROM food WHERE FoodID = ?", "FD001")

	var data entity.Food

	err := row.Scan(&data.FoodID, &data.Name, &data.Serving, &data.Calories)

	log.Println(data)
	// row := config.DB.QueryRow("SELECT * FROM menu WHERE Menu_ID = ?", menuId)

	// var data entities.Menu
	// err := row.Scan(&data.MenuId, &data.MenuName, &data.Price, &data.MenuCategory)
	// log.Println(data.MenuId, data.MenuName, data.Price, data.MenuCategory)

	if err != nil {
		log.Println("fuck")
		log.Println(err)

	}

	return data
}
