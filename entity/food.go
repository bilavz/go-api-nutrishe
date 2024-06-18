package entity

type Food struct {
	FoodID   string `json:"FoodID"`
	Name     string `json:"Name"`
	Serving  int    `json:"Serving"`
	Calories int    `json:"Calories"`
}
