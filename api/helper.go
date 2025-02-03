package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Vulnerability struct {
	ID             string   `json:"id"`
	Severity       string   `json:"severity"`
	CVSS           float64  `json:"cvss"`
	Status         string   `json:"status"`
	PackageName    string   `json:"package_name"`
	CurrentVersion string   `json:"current_version"`
	FixedVersion   string   `json:"fixed_version"`
	Description    string   `json:"description"`
	PublishedDate  string   `json:"published_date"`
	Link           string   `json:"link"`
	RiskFactors    []string `json:"risk_factors"`
}

func processFile(repoURL, fileName string) {
	// scans the repo and stores the file in the database
	rawFileURL := strings.Replace(repoURL, "github.com", "raw.githubusercontent.com", 1)
	rawFileURL = rawFileURL + "/refs/heads/main/" + fileName
	log.Println("Fetching raw file:", rawFileURL)

	var res *http.Response
	var err error

	// try 3 times (1 initial + 1 retry)
	var TRIES = 2
	for i := 0; i < TRIES; i++ {
		res, err = http.Get(rawFileURL)
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

	log.Printf("Response body for file %s: %s", fileName, string(body)[:100])

	var filePayloads []map[string]interface{}
	if err := json.Unmarshal(body, &filePayloads); err != nil {
		log.Default().Printf("Failed to parse file %s: %v", fileName, err)
		return
	}

	// store the file in the database
	for _, payload := range filePayloads {
		log.Printf("Storing payload (JSON): %v", payload)
		storePayload(fileName, payload)
	}
}

func storePayload(fileName string, payload map[string]interface{}) {
	// extract ScanResult from the payload
	scanResults, ok := payload["scanResults"].(map[string]interface{})
	if !ok {
		log.Default().Printf("Failed to extract scan_results from payload")
		return
	}
	log.Printf("\nProcessing scan results: %+v\n", scanResults)

	// extract Vulnerabilities from the scanResults
	vulnerabilities, ok := scanResults["vulnerabilities"].([]interface{})
	if !ok {
		log.Default().Printf("Failed to extract vulnerabilities from scan_results")
		return
	}

	// process each vulnerability and store in the database
	for _, vul := range vulnerabilities {
		vulMap, ok := vul.(map[string]interface{})
		if !ok {
			log.Default().Printf("Failed to convert vulnerability to map")
			continue
		}

		vulnBytes, err := json.Marshal(vulMap)
		if err != nil {
			log.Default().Printf("Failed to marshal vulnerability: %v", err)
			continue
		}

		var vulnerability Vulnerability
		err = json.Unmarshal(vulnBytes, &vulnerability)
		if err != nil {
			log.Default().Printf("Failed to unmarshal vulnerability: %v", err)
			continue
		}

		// Log or process the vulnerability
		log.Printf("Processing vulnerability: %+v", vulnerability)

		// Store the vulnerability in the database
		storeVulnerability(fileName, vulnerability)
	}
}

func storeVulnerability(fileName string, vulnerability Vulnerability) {
	// Store the data in the database
	db, err := sql.Open("sqlite3", "./scans.db")
	if err != nil {
		log.Default().Printf("Failed to open database: %v", err)
		return
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare(`
		INSERT INTO payloads (
			id, severity, cvss, status, package_name, current_version, fixed_version, description, published_date, link, risk_factors, source_file, time_scanned
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		);
	`)
	if err != nil {
		log.Default().Printf("Failed to prepare SQL statement: %v", err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(
		vulnerability.ID,
		vulnerability.Severity,
		vulnerability.CVSS,
		vulnerability.Status,
		vulnerability.PackageName,
		vulnerability.CurrentVersion,
		vulnerability.FixedVersion,
		vulnerability.Description,
		vulnerability.PublishedDate,
		vulnerability.Link,
		strings.Join(vulnerability.RiskFactors, ","),
		fileName,
		time.Now(),
	)
	if err != nil {
		log.Default().Printf("Failed to insert vulnerability into database: %v", err)
	}
}
