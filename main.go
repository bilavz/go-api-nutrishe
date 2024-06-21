package main

import (
	// "go-api-nutrishe/controllers/nabila"

	"log"
	"net/http"
	"nutrishe/controllers/april"

	"nutrishe/controllers/artikel"
	"nutrishe/controllers/mealtrackcontroller"
	"nutrishe/controllers/nabila"
	"nutrishe/controllers/recommendmeals"

	"nutrishe/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Setup database
	err = models.Setup()
	if err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/nutrishe", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/register", nabila.Register)
	mux.HandleFunc("/login", nabila.Login)
	mux.HandleFunc("/dietplan", nabila.CreateDietPlan)
	mux.HandleFunc("/calculate", nabila.CalculateCalories)
	mux.HandleFunc("/calories_goal", nabila.GetCalorieDataHandler)
	mux.HandleFunc("/monthly_calories", nabila.ViewMonthlyCalories)
	mux.HandleFunc("/dailymeal", april.LogMeal)
	mux.HandleFunc("/food", april.GetFoodList)
	mux.HandleFunc("/mealdetail", april.GetMealsByDate)
	mux.HandleFunc("/deletemealdetail", april.DeleteMealDetail)

	mux.HandleFunc("/add_meal", mealtrackcontroller.AddMeal)

	mux.HandleFunc("/recommend_meals", recommendmeals.RecommendMeals)

	mux.HandleFunc("/search_articles", artikel.SearchArticles)

	mux.HandleFunc("/logout", nabila.Logout)
	// Start the HTTP server
	port := ":8081"
	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
