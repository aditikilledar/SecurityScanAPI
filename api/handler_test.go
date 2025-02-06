package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"securityscan/api"
	"securityscan/config"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	os.Setenv("DATABASE_PATH", "./testdb_scans.db")
}

var TEST_DB_PATH = "./testdb_scans.db"

func TestMain(m *testing.M) {

	// log.Print("In run test main!!!")
	// Set the test database path
	os.Setenv("DATABASE_PATH", TEST_DB_PATH)
	// defer os.Unsetenv("DATABASE_PATH")

	// Remove existing test database to start fresh
	os.Remove(TEST_DB_PATH)

	// Initialize the database for testing
	cfg := config.LoadConfig("test")
	cfg.InitDatabaseConfig("test")

	// Run tests
	code := m.Run()

	// Cleanup
	config.CloseDB()
	os.Remove(TEST_DB_PATH)

	os.Exit(code)
}

func TestScanHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    api.ScanRequest
		expectedStatus int
	}{
		{
			name: "Valid request",
			requestBody: api.ScanRequest{
				Repo:      "https://github.com/velancio/vulnerability_scans",
				FileNames: []string{"vulnscan15.json", "vulnscan16.json"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid repo URL",
			requestBody: api.ScanRequest{
				Repo:      "invalid-url",
				FileNames: []string{"file.json"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid request body",
			requestBody:    api.ScanRequest{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty file names",
			requestBody: api.ScanRequest{
				Repo:      "https://github.com/velancio/vulnerability_scans",
				FileNames: []string{},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/scan", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()

			api.ScanHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response api.ScanResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}

				if response.Message == "" {
					t.Error("expected non-empty message in response")
				}

				// Check the contents of the database
				db, err := sql.Open("sqlite3", TEST_DB_PATH)
				if err != nil {
					t.Fatalf("Failed to open database %v", err)
				}

				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM payloads").Scan(&count)
				if err != nil {
					t.Fatalf("Failed to query database: %v", err)
				}
				// hard-coded expected value.
				// 16 is the number of jsons
				if count != 16 {
					t.Errorf("expected %d, got %d", 16, count)
				}
			}
		})
	}
}

func TestQueryHandler(t *testing.T) {

	tests := []struct {
		name           string
		requestBody    api.QueryRequest
		expectedStatus int
	}{
		{
			name: "Valid query",
			requestBody: api.QueryRequest{
				Filters: map[string]string{
					"Severity": "LOW",
				},
			},
			expectedStatus: http.StatusOK,
		},
		// {
		// 	name:           "Invalid request body",
		// 	requestBody:    api.QueryRequest{},
		// 	expectedStatus: http.StatusBadRequest,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}
			log.Print(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()

			api.QueryHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response []api.Vulnerability

				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}

				//check the rows in the database
				db, err := sql.Open("sqlite3", TEST_DB_PATH)
				if err != nil {
					t.Fatalf("Failed to open database %v", err)
				}

				// testing for LOW severity
				var expected_response_len int
				err = db.QueryRow("SELECT COUNT(*) FROM payloads WHERE severity=\"LOW\"").Scan(&expected_response_len)
				if err != nil {
					t.Fatalf("Failed to query database: %v", err)
				}

				if len(response) != expected_response_len {
					t.Errorf("expected %d jsons in response, got %d", expected_response_len, len(response))
				}
			}
		})
	}
}
