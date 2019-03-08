package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/ajiyakin/gohealthz/internal/pkg/handler"
	"github.com/ajiyakin/gohealthz/internal/pkg/storage"
	"github.com/ajiyakin/gohealthz/internal/pkg/updater"
)

var (
	apiJSONRaw []byte
)

func init() {
	var err error
	apiJSONRaw, err = ioutil.ReadFile("api/api.json")
	if err != nil {
		fmt.Printf("unable to load open API specs file from api/api.json")
		os.Exit(1)
	}
}

func main() {
	c, err := parseFlag()
	if err != nil {
		fmt.Printf("invalid flags: %v", err)
		os.Exit(1)
	}

	fmt.Printf("starting service with configurations: %s\n", c.String())

	http.DefaultClient.Timeout = c.httpClientTimeout
	database := storage.NewInMemoryDatabase()

	updater.StartUpdate(database, c.updaterInterval)

	ui := http.FileServer(http.Dir(path.Join("web", "static")))
	http.Handle("/", ui)

	swaggerUI := http.FileServer(http.Dir(path.Join("web", "swagger_ui")))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", swaggerUI))

	http.HandleFunc("/swagger/api.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(apiJSONRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/website", handler.NewWebsiteHandler(database))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
