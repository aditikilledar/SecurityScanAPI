package main

import (
	"log"
	"net/http"
	"os"
	"securityscan/api"
	"securityscan/config"
)

func main() {
	log.SetOutput(os.Stdout)

	cfg := config.LoadConfig("main")

	// init the db
	cfg.InitDatabaseConfig("main")

	defer config.CloseDB()

	// register endpoints
	// TODO: Check if i need to return if endpoint method is not POST
	log.Println("Registering endpoints")

	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received at /scan")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			api.ScanHandler(w, r)
		} else {
			http.Error(w, "Invalid request method, only POST method allowed.", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received at /query")
		if r.Method == http.MethodPost {
			api.QueryHandler(w, r)
		} else {
			http.Error(w, "Invalid request method, only POST method allowed.", http.StatusMethodNotAllowed)
			return
		}
	})

	// manage server
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: nil,
	}

	log.Println("Server started on port 8080")

	// TODO: should this be made HTTPS?
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server Error: %v", err)
	}
}
