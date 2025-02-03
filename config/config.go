package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DatabasePath string
	// Add other configuration fields as needed
}

func LoadConfig() Config {
	// Example of loading configuration from environment variables
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./scans.db"
	}

	return Config{
		DatabasePath: dbPath,
		// Initialize other configuration fields as needed
	}
}

func (c *Config) InitDatabaseConfig() {
	db, err := sql.Open("sqlite3", c.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS payloads (
		id TEXT PRIMARY KEY,
		severity TEXT NOT NULL,
		cvss REAL NOT NULL,
		status TEXT NOT NULL,
		package_name TEXT NOT NULL,
		current_version TEXT NOT NULL,
		fixed_version TEXT NOT NULL,
		description TEXT NOT NULL,
		published_date TEXT NOT NULL,
		link TEXT NOT NULL,
		risk_factors TEXT NOT NULL,
		source_file TEXT NOT NULL,
		time_scanned DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
