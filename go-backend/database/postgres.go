package database

import (
	"database/sql"
	"file-hub-go/config"
	"log"

	_ "github.com/lib/pq" // The PostgreSQL driver
)

var UserDB *sql.DB

// InitUserDB initializes the connection to the PostgreSQL database.
func InitUserDB() {
	var err error
	dbURL := config.AppConfig.DatabaseURL

	UserDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to PostgreSQL database: %v", err)
	}

	if err = UserDB.Ping(); err != nil {
		log.Fatalf("Could not ping PostgreSQL database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database.")
	createUsersTable()
}

// createUsersTable ensures the users table exists.
func createUsersTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`
	if _, err := UserDB.Exec(createTableSQL); err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}
	log.Println("Users table is ready.")
}
