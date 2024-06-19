package models

import (
	"log"
	"nutrishe/entity"
)

func GetDayAndGoal(userID string) entity.DietPlan {
	db := GetDB()
	err := db.QueryRow("SELECT * FROM diet_plan WHERE UserID = ?", userID)

	var temp entity.DietPlan
	err.Scan(&temp.PlanID, &temp.UserID, &temp.StartDate, &temp.EndDate, &temp.CalorieGoal)

	if err != nil {
		log.Println("waduh")
	}

	return temp
}
