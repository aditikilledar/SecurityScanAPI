package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Config struct {
	DatabasePath string
	// Add other configuration fields as needed
}

func LoadConfig(isTest string) Config {
	// Example of loading configuration from environment variables
	var dbPath string
	if isTest == "test" {
		dbPath = "./testdb_scans.db"
	} else {
		dbPath = "./scans.db"
	}

	return Config{
		DatabasePath: dbPath,
		// Initialize other configuration fields as needed
	}
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}
}

func (c *Config) InitDatabaseConfig(isTest string) {
	// var dbPath string
	// if isTest == "test" {
	// 	dbPath = "./testdb_scans.db"
	// 	log.Printf("Using test db %s", dbPath)
	// } else {
	// 	dbPath = "./scans.db"
	// 	log.Printf("Using production db %s", dbPath)
	// }
	var err error
	db, err = sql.Open("sqlite3", c.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS payloads (
        id TEXT NOT NULL,
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
        time_scanned DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (id, source_file)
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func GetDB() *sql.DB {
	return db
}
