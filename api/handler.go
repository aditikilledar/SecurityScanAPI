package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

	// validate if fileNames are json files
	for _, fileName := range fileNames {
		if !strings.HasSuffix(fileName, ".json") {
			http.Error(w, "Invalid file format. Only JSON files are supported at the moment.", http.StatusBadRequest)
			return

		}
	}

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

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	var QueryReq QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&QueryReq); err != nil {
		log.Printf("Error parsing request, invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Processing query request: %+v", QueryReq)

	// process the query and return the response
	vulnerabilities := retrieveFilteredVulnerabilities(QueryReq.Filters)

	if vulnerabilities == nil {
		http.Error(w, "No vulnerabilities found", http.StatusNotFound)
		vulnerabilities = []Vulnerability{}
	}

	log.Printf("Vulnerabilities found in handler: %v", vulnerabilities)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(vulnerabilities); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
