package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"securityscan/api"
)

func TestScanHandler(t *testing.T) {
	tests := []struct {
		name       string
		request    api.ScanRequest
		wantStatus int
	}{
		{
			name: "valid request",
			request: api.ScanRequest{
				Repo:      "https://github.com/velancio/vulnerability_scans",
				FileNames: []string{"vulnscan1011.json", "vulnscan1213.json", "vulnscan456.json"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid repo URL",
			request: api.ScanRequest{
				Repo:      "invalid-url",
				FileNames: []string{"file1.json"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty request",
			request:    api.ScanRequest{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			reqBody, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/scan", bytes.NewBuffer(reqBody))
			rec := httptest.NewRecorder()

			// Call handler
			api.ScanHandler(rec, req)

			// Check status
			if rec.Code != tt.wantStatus {
				t.Errorf("ScanHandler() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}
