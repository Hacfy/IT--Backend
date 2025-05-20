package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

func NewDatabase() *Connection {
	Postgres_uri := os.Getenv("POSTGRES_URI")
	if Postgres_uri == "" {
		log.Fatal("POSTGRES_URI not found")
	}
	DB, err := sql.Open("postgres", Postgres_uri)
	if err != nil {
		log.Fatal(err)
	}
	return &Connection{
		db: DB,
	}
}

func (c *Connection) CheckStatus() {
	if err := c.db.Ping(); err != nil {
		log.Fatal("Bad Database Connection")
	}
	log.Println("Connected to Database successflly")
}

func (c *Connection) Close() {
	if err := c.db.Close(); err != nil {
		log.Fatal("Error while closing the database")
	}
	log.Println("Database Closed")
}
