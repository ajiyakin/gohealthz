package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/ajiyakin/gohealthz/internal/pkg/storage"
	"github.com/google/uuid"
)

type createWebsiteRequest struct {
	URL string `json:"url"`
}

type getWebsitesResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Healty bool   `json:"healty"`
}

// NewWebsiteHandler initilize and get handler for doing website operations
// (POST, GET, DELETE)
func NewWebsiteHandler(database storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createWebsite(w, r, database)
			return
		}
		if r.Method == http.MethodGet {
			getWebsites(w, r, database)
			return
		}
		log.Printf("method %s is not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getWebsites(w http.ResponseWriter, r *http.Request, database storage.Database) {
	websites, err := database.Get()
	if err != nil {
		log.Printf("unable to get list of website from database: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// initilize with make with 0 capacity so if there's no records found,
	// the response body will be [] instead of null
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
}

func createWebsite(w http.ResponseWriter, r *http.Request, database storage.Database) {
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
}
