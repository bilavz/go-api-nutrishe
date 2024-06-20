package models

import (
	"database/sql"
	"nutrishe/entity"
)

func CreateNewMeal(db *sql.DB, data entity.Food) error {
	_, err := db.Exec("INSERT INTO food (FoodID, Name, Serving, Calories, Type) VALUE (?, ?, ?, ?, 'food')", data.FoodID, data.Name, data.Serving, data.Calories)

	if err != nil {
		return err
	}

	return nil
}
