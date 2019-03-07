package main

import (
	"log"
	"net/http"
	"path"
	"time"

	"github.com/ajiyakin/gohealthz/internal/pkg/handler"
	"github.com/ajiyakin/gohealthz/internal/pkg/storage"
	"github.com/ajiyakin/gohealthz/internal/pkg/updater"
)

func main() {
	http.DefaultClient.Timeout = 800 * time.Millisecond
	database := storage.NewInMemoryDatabase()

	updater.StartUpdate(database, 5*time.Minute)

	ui := http.FileServer(http.Dir(path.Join("web", "static")))
	http.Handle("/", ui)
	http.HandleFunc("/website", handler.NewWebsiteHandler(database))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
