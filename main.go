package main

import (
	"log"
	"net/http"

	"securityscan/api"
	"securityscan/config"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library, blank identifier to register the driver
)

func main() {

	config := config.LoadConfig()

	// init the db
	if err := api.InitializeDB(config.DatabasePath); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer api.CloseDB()

	// register endpoints
	// TODO: Check if i need to return if endpoint method is not POST

	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			api.ScanHandler(w, r)
		} else {
			http.Error(w, "Invalid request method, only POST method allowed.", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			api.QueryHandler(w, r)
		} else {
			http.Error(w, "Invalid request method, only POST method allowed.", http.StatusMethodNotAllowed)
			return
		}
	})

	// manage server
	server := &http.Server{}

	log.Println("Server started on port 8080")

	// TODO: should this be made HTTPS?
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server Error: %v", err)
	}
}
