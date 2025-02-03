package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type ScanRequest struct {
	Repo      string   `json:"repo"`
	FileNames []string `json:"files"`
}

type ScanResponse struct {
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}

type QueryRequest struct {
	Filters map[string]string `json:"filters"`
	// focus on severity status
}

func ScanHandler(w http.ResponseWriter, r *http.Request) {
	// Scans JSON files without cloning the git Repo

	// Parse request
	var ScanRequest ScanRequest

	// validate request format
	if err := json.NewDecoder(r.Body).Decode(&ScanRequest); err != nil {
		log.Printf("Error parsing request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	repoURL := ScanRequest.Repo
	fileNames := ScanRequest.FileNames

	// make an HTTP GET request to the repo
	res, err := http.Get(repoURL)

	if (err != nil) || (res.StatusCode != http.StatusOK) {
		// 	http.Error(w, "Invalid repo URL or Failed to Fetch", http.StatusBadRequest)
		// 	return
		http.Error(w, "Invalid repo URL or Failed to Fetch", http.StatusBadRequest)
		return
	}

	log.Printf("Processing scan request - Repo: %s, Files: %v",
		ScanRequest.Repo,
		ScanRequest.FileNames)

	// TODO: Q2 scan files concurrently
	wg := sync.WaitGroup{}
	for _, fileName := range fileNames {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			processFile(repoURL, fileName)
		}(fileName)
	}
	wg.Wait()

	log.Printf("Scan completed successfully at %v", time.Now())
	response := ScanResponse{
		Message:   "Scan completed. Stored files successfully",
		TimeStamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
