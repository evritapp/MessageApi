package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	fmt.Println(connString)
	var err error
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}
	// Test the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error pinging the database: %v", err)
	}

	log.Printf("Connected to the database")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
