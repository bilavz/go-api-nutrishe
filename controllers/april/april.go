package april

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"nutrishe/models"
	"time"
)

type Food struct {
	FoodID        string  `json:"food_id"`
	Name          string  `json:"name"`
	Serving       int     `json:"serving"`
	Calories      int     `json:"calories"`
	Fat           float64 `json:"fat"`
	Carbohydrates float64 `json:"carbohydrates"`
	Protein       float64 `json:"protein"`
	Fiber         *float64 `json:"fiber"`
	Calcium       *int     `json:"calcium"`
	Type          string  `json:"type"`
}

type DailyMeal struct {
	TrackID       string 	`json:"track_id"`
	UserID        string 	`json:"user_id"`
	MealDate      string 	`json:"meal_date"`
	TotalCalories int    	`json:"total_calories"`
}

type MealDetail struct {
	TrackID string `json:"track_id"`
	FoodID  string `json:"food_id"`
}

type DailyMealRequest struct {
    UserID   string `json:"user_id"`
    MealDate string `json:"meal_date"`
    FoodID   string `json:"food_id"`
}

func LogMeal(w http.ResponseWriter, r *http.Request) {
    var mealReq DailyMealRequest
    err := json.NewDecoder(r.Body).Decode(&mealReq)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate and parse the date
    mealDate, err := time.Parse("2006-01-02", mealReq.MealDate)
    if err != nil {
        http.Error(w, "Invalid date format, please use YYYY-MM-DD", http.StatusBadRequest)
        return
    }

    db := models.GetDB()
    if db == nil {
        http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
        return
    }

    // Start transaction to ensure atomicity
    tx, err := db.Begin()
    if err != nil {
        log.Printf("Failed to start transaction: %v", err)
        http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
        return
    }
    defer func() {
        if err != nil {
            log.Printf("Transaction failed, rolling back: %v", err)
            tx.Rollback()
            http.Error(w, "Transaction failed", http.StatusInternalServerError)
        }
    }()

    // Check if daily meal already exists for the user and meal date
    var trackID string
    var totalCalories int
    query := "SELECT TrackID, TotalCalories FROM daily_meal WHERE UserID = ? AND MealDate = ?"
    err = tx.QueryRow(query, mealReq.UserID, mealDate).Scan(&trackID, &totalCalories)
    if err != nil && err != sql.ErrNoRows {
        log.Printf("Failed to check existing daily meal: %v", err)
        http.Error(w, "Failed to check existing daily meal", http.StatusInternalServerError)
        return
    }

    if trackID == "" {
        // Generate new TrackID
        trackID, err = models.GenerateSequentialTrackID()
        if err != nil {
            log.Printf("Failed to generate TrackID: %v", err)
            http.Error(w, "Failed to generate TrackID", http.StatusInternalServerError)
            return
        }

        // Insert new daily meal
        _, err = tx.Exec("INSERT INTO daily_meal (TrackID, UserID, MealDate, TotalCalories) VALUES (?, ?, ?, 0)", trackID, mealReq.UserID, mealDate)
        if err != nil {
            log.Printf("Failed to create daily meal: %v", err)
            http.Error(w, "Failed to create daily meal", http.StatusInternalServerError)
            return
        }
    }

    // Check if food_id exists in the food table
    var foodExists bool
    err = tx.QueryRow("SELECT COUNT(*) > 0 FROM food WHERE FoodID = ?", mealReq.FoodID).Scan(&foodExists)
    if err != nil || !foodExists {
        log.Printf("Invalid food_id: %v", mealReq.FoodID)
        http.Error(w, "Invalid food_id", http.StatusBadRequest)
        return
    }

    // Insert meal detail
    _, err = tx.Exec("INSERT INTO meal_detail (TrackID, FoodID) VALUES (?, ?)", trackID, mealReq.FoodID)
    if err != nil {
        log.Printf("Failed to add meal detail: %v", err)
        http.Error(w, "Failed to add meal detail", http.StatusInternalServerError)
        return
    }

    // Calculate updated total calories
    var newTotalCalories int
    err = tx.QueryRow("SELECT SUM(f.Calories) FROM meal_detail md JOIN food f ON md.FoodID = f.FoodID WHERE md.TrackID = ?", trackID).Scan(&newTotalCalories)
    if err != nil {
        log.Printf("Failed to calculate total calories: %v", err)
        http.Error(w, "Failed to calculate total calories", http.StatusInternalServerError)
        return
    }

    // Update total calories in daily_meal table
    _, err = tx.Exec("UPDATE daily_meal SET TotalCalories = ? WHERE TrackID = ?", newTotalCalories, trackID)
    if err != nil {
        log.Printf("Failed to update total calories: %v", err)
        http.Error(w, "Failed to update total calories", http.StatusInternalServerError)
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        log.Printf("Failed to commit transaction: %v", err)
        http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Meal logged successfully"})
}

func GetFoodList(w http.ResponseWriter, r *http.Request) {
	db := models.GetDB()
	rows, err := db.Query("SELECT FoodID, Name, Serving, Calories, Fat, Carbohydrates, Protein, Fiber, Calcium, Type FROM food")
	if err != nil {
		http.Error(w, "Failed to retrieve food list: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

    var foodsByType = make(map[string][]Food)

	for rows.Next() {
        var food Food
        var calcium sql.NullInt64
        var fiber sql.NullFloat64
    
        if err := rows.Scan(&food.FoodID, &food.Name, &food.Serving, &food.Calories, &food.Fat, &food.Carbohydrates, &food.Protein, &fiber, &calcium, &food.Type); err != nil {
            http.Error(w, "Failed to scan food item: "+err.Error(), http.StatusInternalServerError)
            return
        }
    
        if fiber.Valid {
            food.Fiber = &fiber.Float64
        } else {
            food.Fiber = nil
        }
    
        if calcium.Valid {
            calciumValue := int(calcium.Int64)
            food.Calcium = &calciumValue
        } else {
            food.Calcium = nil
        }
    
        foodsByType[food.Type] = append(foodsByType[food.Type], food)
    }
    
    

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foodsByType)
}

func GetMealsByDate(w http.ResponseWriter, r *http.Request) {
    var requestBody struct {
        UserID   string `json:"user_id"`
        MealDate string `json:"meal_date"`
    }

    // Decode the JSON request body
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    userID := requestBody.UserID
    mealDateStr := requestBody.MealDate

    if userID == "" || mealDateStr == "" {
        http.Error(w, "Missing user_id or meal_date", http.StatusBadRequest)
        return
    }

    // Validate and parse the date
    mealDate, err := time.Parse("2006-01-02", mealDateStr)
    if err != nil {
        http.Error(w, "Invalid date format, please use YYYY-MM-DD", http.StatusBadRequest)
        return
    }

    db := models.GetDB()
    if db == nil {
        http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
        return
    }

    var trackID string
    var totalCalories int
    err = db.QueryRow("SELECT TrackID, TotalCalories FROM daily_meal WHERE UserID = ? AND MealDate = ?", userID, mealDate).Scan(&trackID, &totalCalories)
    if err != nil {
        log.Printf("Failed to retrieve track ID: %v", err)
        http.Error(w, "Failed to retrieve track ID", http.StatusInternalServerError)
        return
    }

    rows, err := db.Query("SELECT f.FoodID, f.Name, f.Serving, f.Calories, f.Fat, f.Carbohydrates, f.Protein, f.Fiber, f.Calcium, f.Type FROM meal_detail md JOIN food f ON md.FoodID = f.FoodID WHERE md.TrackID = ?", trackID)
    if err != nil {
        log.Printf("Failed to retrieve meals: %v", err)
        http.Error(w, "Failed to retrieve meals", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var meals []Food
    for rows.Next() {
        var food Food
        var calcium sql.NullInt64
        var fiber sql.NullFloat64

        if err := rows.Scan(&food.FoodID, &food.Name, &food.Serving, &food.Calories, &food.Fat, &food.Carbohydrates, &food.Protein, &fiber, &calcium, &food.Type); err != nil {
            log.Printf("Failed to scan meal item: %v", err)
            http.Error(w, "Failed to scan meal item", http.StatusInternalServerError)
            return
        }

        if fiber.Valid {
            food.Fiber = &fiber.Float64
        } else {
            food.Fiber = nil
        }

        if calcium.Valid {
            calciumValue := int(calcium.Int64)
            food.Calcium = &calciumValue
        } else {
            food.Calcium = nil
        }

        meals = append(meals, food)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating over rows: %v", err)
        http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
        return
    }

    // Jika tidak ada makanan yang ditemukan, kirim respons kosong (tidak ada error)
    if len(meals) == 0 {
        meals = []Food{}
    }

    response := struct {
        TotalCalories int    `json:"total_calories"`
        Meals         []Food `json:"meals"`
    }{
        TotalCalories: totalCalories,
        Meals:         meals,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

