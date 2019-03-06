package main

import (
	"log"
	"net/http"
	"path"

	"github.com/ajiyakin/gohealthz/internal/pkg/handler"
	"github.com/ajiyakin/gohealthz/internal/pkg/storage"
)

func main() {
	database := storage.NewInMemoryDatabase()

	ui := http.FileServer(http.Dir(path.Join("web", "static")))
	http.Handle("/", ui)
	http.HandleFunc("/website", handler.NewWebsiteHandler(database))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
