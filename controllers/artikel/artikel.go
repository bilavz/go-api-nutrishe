package artikel

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	apiKey = "AIzaSyBG6ee5dKpxZ-0kHbAyzpG8WslqSzdWF_c"
	cx     = "84161d7b25ebd4a92"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

type SearchResult struct {
	Kind string `json:"kind"`
	Url  struct {
		Type     string `json:"type"`
		Template string `json:"template"`
	} `json:"url"`
	Queries struct {
		Request []struct {
			Title        string `json:"title"`
			TotalResults string `json:"totalResults"`
			SearchTerms  string `json:"searchTerms"`
			Count        int    `json:"count"`
			StartIndex   int    `json:"startIndex"`
		} `json:"request"`
	} `json:"queries"`
	Items []struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
	} `json:"items"`
}

func SearchArticles(w http.ResponseWriter, r *http.Request) {
	searchQuery := "Healthy food recommendation" //search query nya
	maxResults := 5                              //maxnya mau brp artikel
	searchResults, err := search(searchQuery, maxResults)
	if err != nil {
		log.Fatalf("Error performing search: %v", err)
	}

	log.Println("Search Results:")
	for _, item := range searchResults.Items {
		log.Printf("Title: %s\n", item.Title)
		log.Printf("Link: %s\n", item.Link)
		log.Printf("Snippet: %s\n\n", item.Snippet)
	}

	respondJSON(w, http.StatusOK, searchResults.Items)
}

func search(query string, maxResults int) (*SearchResult, error) {
	baseURL := "https://www.googleapis.com/customsearch/v1"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("key", apiKey)
	q.Set("cx", cx)
	q.Set("q", query)
	q.Set("num", fmt.Sprintf("%d", maxResults))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
