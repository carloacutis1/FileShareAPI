package database

// connects to the database and provides a global DB variable

// to create database:
// brew services start postgresql@18
// createdb database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Init opens the pool and pings Postgres. Call once from main.
func Init() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default local Postgres on port 5433 — change user/password/db to match your machine
		dsn = "postgres://postgres:postgres@localhost:5433/database?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("sql.Open:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("db.Ping:", err)
	}

	DB = db
}
