package main

import (
	// "go-api-nutrishe/controllers/nabila"

	"log"
	"net/http"
	"nutrishe/controllers/nabila"
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
	mux.HandleFunc("/empowher", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/register", nabila.Register)
	mux.HandleFunc("/login", nabila.Login)
	mux.HandleFunc("/dietplan", nabila.CreateDietPlan)
	mux.HandleFunc("/calculate_calories", nabila.CalculateCalories)
	mux.HandleFunc("/calories_goal", nabila.ViewCaloriesGoal)
	mux.HandleFunc("/monthly_calories", nabila.ViewMonthlyCalories)

	// Start the HTTP server
	port := ":8081"
	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// func loadEnv() error {
// 	// Simulasi membaca environment variable, sebaiknya gunakan library seperti godotenv untuk membaca file .env
// 	jwtKey := os.Getenv("JWT_KEY")
// 	if jwtKey == "" {
// 		return fmt.Errorf("JWT_KEY environment variable not set")
// 	}

// 	return nil
// }
