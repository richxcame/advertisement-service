package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig holds the database configuration parameters
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

// Connect connects to the database and returns a *sql.DB instance
func Connect(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate creates the advertisements table if it does not exist
func Migrate(db *sql.DB, query string) {
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
