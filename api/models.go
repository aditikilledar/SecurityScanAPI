package api

import (
	"errors"
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

// TODO: check if embedded struct works
type Payload struct {
	Vulnerability        // Embedded struct
	SourceFile    string `json:"source_file"`
	TimeScanned   string `json:"time_scanned"`
}

func (p *Payload) ToVulnerability() Vulnerability {
	return p.Vulnerability
}

// Validate the Payload struct
func (p *Payload) Validate() error {
	if p.ID == "" || p.Severity == "" || p.CVSS == 0 || p.Status == "" || p.PackageName == "" || p.CurrentVersion == "" || p.FixedVersion == "" || p.Description == "" || p.PublishedDate == "" || p.Link == "" {
		return errors.New("missing required fields")
	}
	return nil
}

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

type QueryResponse struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}
