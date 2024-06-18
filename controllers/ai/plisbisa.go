package ai

import (
	"context"
	"os"
	"strconv"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	mealplan "nutrishe/models"

	"log"
)

var aiAnswered = false

var aiSaid string

func TestAI() {
	log.Println("tes ai")

	userID := "US001" //ini ceritanya buat ngetes

	// if r.Method == "GET" {
	// 	log.Println("get")

	// 	if !aiAnswered {
	// 		temp, _ := template.ParseFiles("views/index.html")
	// 		temp.Execute(w, nil)
	// 	} else {
	// 		data := map[string]any{
	// 			"data": aiSaid,
	// 		}

	// 		temp, _ := template.ParseFiles("views/index.html")
	// 		temp.Execute(w, data)
	// 	}
	// }

	// if r.Method == "POST" {
	// log.Println("post")

	os.Setenv("API_KEY", "AIzaSyBe8wUJ5RsbBqPSJMtFKV5BnkcrF93u_8o")

	// r.ParseForm()
	// action := r.Form.Get("action")
	// log.Println(action)

	action := " "
	levelTest := "beginner"

	var prompt string

	if action == "recommend" {
		prompt = "Generate healthy foods recommendation for people with an experience level of "

		// temp := r.Form.Get("level")

		// prompt += temp
		prompt += levelTest

	} else if action == "dietPlan" {
		prompt = "Generate a meal plan for  "
		//days

		temp := "2" //buat tes, aku gatau carane dpt hari diantara 2 tanggal
		prompt += temp

		prompt += " days, with "
		//caories

		mealplan := mealplan.GetDayAndGoal(userID)

		calories := strconv.Itoa(mealplan.CalorieGoal)

		prompt += calories

		prompt += " each day"
	}

	// log.Println("Prompt", prompt)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, _ := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))

	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, _ := model.GenerateContent(ctx, genai.Text("Generate foods with "))

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

	// prompt: Generate healthy foods recommendation for people with an experience level of beginner, specified to Indonesian dishes, display calories for each food item
	// prompt: Generate a meal plan for ? day, ? calories each day, with calories for each day. Specific to Indonesian dishes.

	// data := map[string]any{
	// 	"data": aiResponse,
	// }

	// temp, _ := template.ParseFiles("views/index.html")
	// temp.Execute(w, data)
	// }

	// temp, _ := template.ParseFiles("views/reservation/menu.html")
	// temp.Execute(w, nil)
}
