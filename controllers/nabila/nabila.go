package nabila

import (
	"encoding/json"
	"net/http"
	"nutrishe/models"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"fmt"
)

// Credentials represents the credentials required for login and register
type Credentials struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birthdate"`
	Height    float32   `json:"height"`
	Weight    float32   `json:"weight"`
}

type RegisterCredentials struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Birthdate string  `json:"birthdate"`
	Height    float32 `json:"height"`
	Weight    float32 `json:"weight"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))

// DailyLogRequest represents the request payload for the daily log
type DailyLogRequest struct {
	UserID     string `json:"user_id"`
	SymptomsID string `json:"symptoms_id"`
}

// CreateDietPlanRequest represents the request payload for creating a diet plan
type CreateDietPlanRequest struct {
	UserID      string `json:"user_id"`
	Description string `json:"description"`
	CalorieGoal int    `json:"calorie_goal"`
}

// CalorieRequest represents the request payload for calculating calories
type CalorieRequest struct {
	UserID string `json:"user_id"`
	Date   string `json:"date"`
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var creds RegisterCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse birthdate
	birthdate, err := time.Parse("2006-01-02", creds.Birthdate)
	if err != nil {
		http.Error(w, "Invalid birthdate format", http.StatusBadRequest)
		return
	}

	userID := generateID()

	db := models.GetDB()
	err = models.CreateUser(db, userID, creds.Name, creds.Username, creds.Email, creds.Password, birthdate, creds.Height, creds.Weight)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	user, err := models.GetUserByEmail(db, creds.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Passwords), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.UserID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Write([]byte(tokenString))
}

// generateID generates a unique ID for the user

func generateID() string {
	bytes := make([]byte, 3) // 3 bytes menghasilkan angka 6 karakter
	rand.Read(bytes)
	return fmt.Sprintf("US%03d", bytes[0]<<8|bytes[1])
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// CreateDietPlan handles creating a new diet plan
func CreateDietPlan(w http.ResponseWriter, r *http.Request) {
	var req CreateDietPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Bad request"})
		return
	}

	if req.UserID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "UserID is required"})
		return
	}

	dietPlan := models.DietPlan{
		PlanID:      generateID(), // Assume a function to generate unique IDs
		UserID:      req.UserID,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0), // Assuming a default one month duration
		CalorieGoal: req.CalorieGoal,
	}

	db := models.GetDB()
	query := "INSERT INTO diet_plan (PlanID, UserID, StartDate, EndDate, CalorieGoal) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, dietPlan.PlanID, dietPlan.UserID, dietPlan.StartDate, dietPlan.EndDate, dietPlan.CalorieGoal)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error creating diet plan", "error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, dietPlan)
}

// CalculateCalories handles calculating the calories for a given date
func CalculateCalories(w http.ResponseWriter, r *http.Request) {
	var req CalorieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Bad request"})
		return
	}

	if req.UserID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "UserID is required"})
		return
	}

	var totalCalories int
	db := models.GetDB()
	query := "SELECT SUM(total_calories) FROM daily_meal WHERE user_id = ? AND meal_date = ?"
	err := db.QueryRow(query, req.UserID, req.Date).Scan(&totalCalories)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error calculating calories", "error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"total_calories": totalCalories})
}

// ViewCaloriesGoal handles viewing the calorie goal for a user
func ViewCaloriesGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "UserID is required"})
		return
	}

	var calorieGoal int
	db := models.GetDB()
	query := "SELECT CalorieGoal FROM dietplan WHERE UserID = ? ORDER BY EndDate DESC LIMIT 1"
	err := db.QueryRow(query, userID).Scan(&calorieGoal)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error retrieving diet plan", "error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"calorie_goal": calorieGoal})
}

// ViewMonthlyCalories handles viewing the total calories for the current month
func ViewMonthlyCalories(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "UserID is required"})
		return
	}

	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var totalCalories int
	db := models.GetDB()
	query := "SELECT SUM(total_calories) FROM daily_meal WHERE user_id = ? AND meal_date BETWEEN ? AND ?"
	err := db.QueryRow(query, userID, startOfMonth, endOfMonth).Scan(&totalCalories)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error calculating monthly calories", "error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"total_calories": totalCalories})
}
