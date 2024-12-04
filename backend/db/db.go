package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/UpsDev42069/BM_Search_Engine/backend/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDBenv() {
	config.LoadEnv()
}

// ConnectDB returns a new connection to the database.
func ConnectDB(initMode bool) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if config.DBDriver == "sqlite3" {
		// For SQLite in-memory database
		db, err = sql.Open("sqlite3", ":memory:")
	} else {
		// For PostgreSQL
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
		db, err = sql.Open("postgres", psqlInfo)
	}
	if err != nil {
		return nil, err
	}

	if !initMode {
		if err := CheckDBExists(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// CheckDBExists checks if the database exists (only applicable for PostgreSQL).
func CheckDBExists(db *sql.DB) error {
	if config.DBDriver == "sqlite3" {
		// SQLite does not support databases in the same way; skip existence check
		return nil
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)"
	err := db.QueryRow(query, config.DBName).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("database %s does not exist", config.DBName)
	}

	return nil
}

// QueryDB queries the database and returns a list of maps.
func QueryDB(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Initialize an empty slice to hold the query results
	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		// Creates a slice to hold the column values
		values := make([]interface{}, len(columns))
		// Creates a slice of pointers to the column values
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Creates a map to hold the column name-value pairs
		result := make(map[string]interface{})
		for i, col := range columns {
			result[col] = values[i]
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetUserID looks up the id for a username.
func GetUserID(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RunMigrations runs the database migrations.
func RunMigrations(db *sql.DB) {
	var dbURL string
	if config.DBDriver == "sqlite3" {
		dbURL = "sqlite3://file::memory:?cache=shared" // SQLite in-memory DB
	} else {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	}

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrated successfully")
}
