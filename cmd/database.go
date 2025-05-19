package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Connection struct {
	db *sql.DB
}

func NewDatabase() *Connection {
	DB, err := sql.Open("postgres", "postgresql://it_inventory_mic8_user:DkyTzw7x3PtUTzdDFuNLoSqEBXRQgJG5@dpg-d0kscfruibrs739q4pl0-a.singapore-postgres.render.com/it_inventory_mic8")
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
