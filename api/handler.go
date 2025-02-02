package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

//    "repo": <repo root>,   "files": [<filename1>, <filename2>,…]}

type ScanRequest struct {
	Repo      string   `json:"repo"`
	FileNames []string `json:"files"`
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

	// TODO: Q2 scan files concurrently
	wg := sync.WaitGroup{}
	for _, fileName := range fileNames {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()

			storeFile(repoURL, fileName)
		}(fileName)
	}
	wg.Wait()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Scan completed successfully"))
}
