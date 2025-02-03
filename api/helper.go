package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func storeFile(repoURL, fileName string) {
	// scans the repo and stores the file in the database
	fileURL := repoURL + "/" + fileName
	var res *http.Response
	var err error

	// try 3 times (1 initial + 1 retry)
	var TRIES = 2
	for i := 0; i < TRIES; i++ {
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

	// read the response body
	body, err := io.ReadAll(res.Body)
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
	payloadBytes, err := json.Marshal(payload) // Store the marshalled payload in a new variable
	if err != nil {
		log.Default().Printf("Failed to marshal payload: %v", err)
		return
	}
	// TODO: Design Schema for :
	// source file, scan time and payload
	log.Println("STORED @", time.Now(), fileName, " IN DB (Fake)", string(payloadBytes)[:100])

	// TODO: Logic to store into database
	// _, err := db.Exec(SCAN_INSERT_QUERY, fileName, payloadBytes)
}
