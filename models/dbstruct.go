package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Cycle represents the cycle table
type Cycle struct {
	CycleID       string    `json:"cycle_id" gorm:"type:char(5);primaryKey"`
	UserID        string    `json:"user_id" gorm:"type:char(5)"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CycleDuration int       `json:"cycle_duration"`
}

// DailyLog represents the daily_log table
type DailyLog struct {
	CycleID    string `json:"cycle_id" gorm:"type:char(5);primaryKey"`
	SymptomsID string `json:"symptoms_id" gorm:"type:char(5);primaryKey"`
}

// DailyMeal represents the daily_meal table
type DailyMeal struct {
	TrackID       string    `json:"track_id" gorm:"type:char(5);primaryKey"`
	UserID        string    `json:"user_id" gorm:"type:char(5)"`
	MealDate      time.Time `json:"meal_date"`
	TotalCalories int       `json:"total_calories"`
}

// DietPlan represents the diet_plan table
type DietPlan struct {
	PlanID      string    `json:"plan_id" gorm:"type:char(5);primaryKey"`
	UserID      string    `json:"user_id" gorm:"type:char(5)"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CalorieGoal int       `json:"calorie_goal"`
}

// Food represents the food table
type Food struct {
	FoodID        string  `json:"food_id" gorm:"type:char(5);primaryKey"`
	Name          string  `json:"name" gorm:"type:varchar(255)"`
	Serving       int     `json:"serving"`
	Calories      int     `json:"calories"`
	Fat           float32 `json:"fat"`
	Carbohydrates float32 `json:"carbohydrates"`
	Protein       float32 `json:"protein"`
	Fiber         float32 `json:"fiber"`
	Calcium       int     `json:"calcium"`
}

// MealDetail represents the meal_detail table
type MealDetail struct {
	TrackID  string    `json:"track_id" gorm:"type:char(5);primaryKey"`
	FoodID   string    `json:"food_id" gorm:"type:char(5);primaryKey"`
	MealTime time.Time `json:"meal_time"`
}

// Migration represents the migrations table
type Migration struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Migrations string `json:"migrations" gorm:"type:varchar(255)"`
	Batch      int    `json:"batch"`
}

// SymptomsType represents the symptoms_type table
type SymptomsType struct {
	SymptomsID   string `json:"symptoms_id" gorm:"type:char(5);primaryKey"`
	Category     string `json:"category" gorm:"type:varchar(55)"`
	SymptomsName string `json:"symptoms_name" gorm:"type:varchar(255)"`
}

// User represents the users table
type User struct {
	UserID    string    `json:"user_id" gorm:"type:char(5);primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	Username  string    `json:"username" gorm:"type:varchar(55)"`
	Email     string    `json:"email" gorm:"type:varchar(55)"`
	Passwords string    `json:"passwords" gorm:"type:varchar(55)"`
	Birthdate time.Time `json:"birthdate"`
	Height    float32   `json:"height"`
	Weight    float32   `json:"weight"`
}

func CreateUser(db *sql.DB, userID, name, username, email, password string, birthdate time.Time, height, weight float32) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (UserID, name, username, email, passwords, birthdate, height, weight) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, userID, name, username, email, hashedPassword, birthdate, height, weight)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := "SELECT UserID, Name, Username, Email, Passwords, Birthdate, Height, Weight FROM users WHERE Email = ?"
	row := db.QueryRow(query, email)

	var user User
	var birthdateStr string
	err := row.Scan(&user.UserID, &user.Name, &user.Username, &user.Email, &user.Passwords, &birthdateStr, &user.Height, &user.Weight)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Pengguna tidak ditemukan
		}
		return nil, err // Kesalahan lainnya
	}

	// Konversi string ke time.Time
	user.Birthdate, err = time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
