package recommendmeals

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"nutrishe/entity"

	"log"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func RecommendMeals(w http.ResponseWriter, r *http.Request) {
	log.Println("rekomen ai")

	var data entity.AIPrompt
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prompt: Generate a meal plan for ? day, ? calories each day, with calories for each meal. Specific to ? dishes.

	//making prompt
	prompt := "Generate a meal plan for  " + data.Days + " days, " + data.Calories +
		" calories each day, with calories for each meal. Specific to " + data.Cuisine + " cuisines."

	log.Print("prompt: ", prompt)

	ctx := context.Background()
	client, _ := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("Genai_API_KEY")))

	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, _ := model.GenerateContent(ctx, genai.Text(prompt))

	var aiResponse genai.Part
	if resp != nil {
		candidates := resp.Candidates

		for _, candidate := range candidates {
			content := candidate.Content
			if content != nil {
				aiResponse = content.Parts[0]
				log.Print("AI response:  ", aiResponse)
			}
		}
	}

	respondJSON(w, http.StatusOK, aiResponse)
}
