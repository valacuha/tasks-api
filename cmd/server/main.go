package main

import (
	"log"
	"net/http"
	"tasks-api/internal/handlers"
	middleware "tasks-api/internal/http"
	"tasks-api/internal/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	h := handlers.New(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.TasksCollection)
	mux.HandleFunc("/tasks/", h.TaskItem)
	mux.HandleFunc("/health", h.Health)

	loggedMux := middleware.LoggingMiddleware(mux)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatal(err)
	}
}
