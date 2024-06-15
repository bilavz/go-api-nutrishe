package nabila

import (
	"encoding/json"
	"net/http"
	"nutrishe/models"
	"time"

	"fmt"
)

// DailyLogRequest represents the request payload for the daily log
type DailyLogRequest struct {
	UserID     string `json:"user_id"`
	SymptomsID string `json:"symptoms_id"`
}

// Checkin handles the daily log request
func Checkin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req DailyLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "CycleID is required", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	query := "INSERT INTO daily_log (CycleID, SymptompsID, check_in_date) VALUES (?, ?, ?)"
	_, err := db.Exec(query, req.UserID, req.SymptomsID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting check-in data: %v", err), http.StatusInternalServerError)
		return
	}

	// dailyLog := models.DailyLog{
	// 	CycleID:    generateID(), // Assume a function to generate unique IDs
	// 	SymptomsID: req.SymptomsID,
	// }

	// if err := models.DB.Create(&dailyLog).Error; err != nil {
	// 	http.Error(w, "Failed to log data", http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(dailyLog)
	w.Write([]byte("Check-in successful"))
}

// generateID generates a unique ID for new entries
func generateID() string {
	// Implement your ID generation logic here
	return "some_unique_id"
}
