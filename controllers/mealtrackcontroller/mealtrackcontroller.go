package mealtrackcontroller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nutrishe/entity"
	"nutrishe/models"
)

func generateFoodID() string {
	bytes := make([]byte, 3) // 3 bytes menghasilkan angka 6 karakter
	rand.Read(bytes)
	return fmt.Sprintf("FD%03d", bytes[0]<<8|bytes[1])
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func AddMeal(w http.ResponseWriter, r *http.Request) {
	log.Println("Add a meal")
	var data entity.Food
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	foodID := generateFoodID()

	data.FoodID = foodID

	log.Println(data)

	db := models.GetDB()
	err = models.CreateNewMeal(db, data)
	if err != nil {
		http.Error(w, "Failed to create new meal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
