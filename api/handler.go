package api

import (
	"encoding/json"
	"io"
	"log"
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

func storeFile(repoURL, fileName string) {
	// scans the repo and stores the file in the database
	fileURL := repoURL + "/" + fileName
	var res *http.Response
	var err error

	// try 3 times
	for i := 0; i < 3; i++ {
		res, err = http.Get(fileURL)
		if err == nil {
			break
		}
	}

	if err != nil || res.StatusCode != http.StatusOK {
		log.Default().Printf("Failed to fetch file %s: %v", fileName, err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body) // Read the body into a byte slice
	if err != nil {
		log.Default().Printf("Failed to read response body for file %s: %v", fileName, err)
		return
	}

	var filePayloads []map[string]interface{}
	if err := json.Unmarshal(body, &filePayloads); err != nil {
		log.Default().Printf("Failed to parse file %s: %v", fileName, err)
	}

	// store the file in the database
	for _, payload := range filePayloads {
		storePayload(fileName, payload)
	}
}

func storePayload(fileName string, payload map[string]interface{}) {
	payload, err := json.Marshal(payload)
	if err != nil {
		log.Default().Printf("Failed to marshal payload: %v", err)
		return
	}

	// TODO: Logic to store into database
	_, err := db.Exec(SCAN_INSERT_QUERY, fileName, payload)
}
