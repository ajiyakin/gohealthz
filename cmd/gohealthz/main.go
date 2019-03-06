package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/ajiyakin/gohealthz/internal/pkg/storage"

	"github.com/google/uuid"
)

var (
	database = storage.NewInMemoryDatabase()
)

func main() {
	ui := http.FileServer(http.Dir(path.Join("web", "static")))
	http.Handle("/", ui)
	http.HandleFunc("/website", createWebsiteHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type createWebsiteRequest struct {
	URL string `json:"url"`
}

type getWebsitesResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Healty bool   `json:"healty"`
}

func createWebsiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var requestBody createWebsiteRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Printf("unable to decode request body: %v", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		id, err := uuid.NewUUID()
		if err != nil {
			log.Printf("unable to generate new UUID: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_, err = url.ParseRequestURI(requestBody.URL)
		if err != nil {
			log.Printf("unable to parse URL: %v with URL input: %s", err, requestBody.URL)
			http.Error(w, "invalid URL. URL must be in form of absolute URL", http.StatusBadRequest)
			return
		}
		// TODO Set timeout to 800ms
		var healthiness bool
		response, err := http.Get(requestBody.URL)
		if err != nil {
			log.Printf("url %s is not healthy: %v", requestBody.URL, err)
		}
		if err == nil && response.StatusCode == http.StatusOK {
			healthiness = true
		}
		err = database.Save(storage.Website{
			ID:      id.String(),
			URL:     requestBody.URL,
			Healthy: healthiness,
		})
		if err != nil {
			log.Printf("unable to save to database: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Print("successfully store website to database")
		w.WriteHeader(http.StatusCreated)
		return
	}
	if r.Method == http.MethodGet {
		websites, err := database.Get()
		if err != nil {
			log.Printf("unable to get list of website from database: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		responseBody := make([]getWebsitesResponse, 0)
		for _, website := range websites {
			responseBody = append(responseBody, getWebsitesResponse{
				ID:     website.ID,
				URL:    website.URL,
				Healty: website.Healthy,
			})
		}
		w.Header().Add("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(&responseBody); err != nil {
			log.Printf("unable to encode records to response writter: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Print("successfully retrieve website records")
		return
	}
	log.Printf("method %s is not allowed", r.Method)
	w.WriteHeader(http.StatusMethodNotAllowed)
}
