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

func getFilteredPayloadsBySeverity(severity string) ([]Payload, error) {
	db, err := sql.Open("sqlite3", "./scans.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	defer db.Close()

	FILTER_SQL_STMT := `SELECT id, severity, cvss, status, package_name, current_version, fixed_version, description, published_date, link, risk_factors, source_file, time_scanned FROM payloads 
	WHERE severity = ? 
	AND severity != ''`

	rows, err := db.Query(FILTER_SQL_STMT, severity)
	if err != nil {
		log.Printf("Error querying the database: %v", err)
		return nil, err
	}

	var payloads []Payload
	for rows.Next() {
		var payload Payload
		var risk_factors string
		if err := rows.Scan(
			&payload.ID,
			&payload.Severity,
			&payload.CVSS,
			&payload.Status,
			&payload.PackageName,
			&payload.CurrentVersion,
			&payload.FixedVersion,
			&payload.Description,
			&payload.PublishedDate,
			&payload.Link,
			&risk_factors,
			&payload.SourceFile,
			&payload.TimeScanned,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		payload.RiskFactors = parseRiskFactors(risk_factors)
		log.Printf("Retreieved payload >>> %v", payload)
		payloads = append(payloads, payload)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return payloads, nil
}

func parseRiskFactors(riskFactors string) []string {
	return strings.Split(riskFactors, ",")
}

func retrieveFilteredVulnerabilities(Filters map[string]string) []Vulnerability {
	// takes a set of filters and returns the filtered payloads as an sql rows pointer

	// currently handles only severity filter
	// TODO: Handle other filters
	severity := Filters["severity"]
	log.Printf("Filtering query request for severity: %v", severity)

	payloads, err := getFilteredPayloadsBySeverity(severity)

	if err != nil {
		log.Printf("Error retrieving filtered payloads: %v", err)
	}

	log.Printf("Retrieved filtered payloads: %v", len(payloads))

	// convert payloads to vulnerabilities
	var vulnerabilities []Vulnerability

	for _, payload := range payloads {
		if payload.Severity == "" {
			continue
		}
		log.Printf("Appending vul: %v", payload.Vulnerability)
		vulnerabilities = append(vulnerabilities, payload.Vulnerability)
	}

	return vulnerabilities
}
