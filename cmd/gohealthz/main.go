package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/ajiyakin/gohealthz/internal/pkg/handler"
	"github.com/ajiyakin/gohealthz/internal/pkg/storage"
	"github.com/ajiyakin/gohealthz/internal/pkg/updater"
)

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
	http.HandleFunc("/website", handler.NewWebsiteHandler(database))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
