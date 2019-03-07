package updater

import (
	"log"
	"net/http"
	"time"

	"github.com/ajiyakin/gohealthz/internal/pkg/storage"
)

// StartUpdate starts (run) updater on the background and will update
// website healthiness in a given interval
func StartUpdate(database storage.Database, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func(database storage.Database) {
		for range ticker.C {
			updateHealthiness(database)
		}
	}(database)
}

func updateHealthiness(database storage.Database) {
	websites, err := database.Get()
	if err != nil {
		log.Printf("unable to get list of websites to update: %v", err)
		return
	}
	for _, website := range websites {
		response, err := http.Get(website.URL)
		if err != nil {
			log.Printf("unable to get request for URL: %s. error: %v", website.URL, err)
			website.Healthy = false
			continue
		}
		if response.StatusCode != http.StatusOK {
			log.Printf("website with URL: %s is not healthy. response code: %d", website.URL, response.StatusCode)
			website.Healthy = false
		}
		if err = database.Save(website); err != nil {
			log.Printf("unable to save (update) to database: %v", err)
		}
	}
}
