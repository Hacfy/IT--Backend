package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	addr := flag.String("addr", ":8080", "Address at which the server will run")
	flag.Parse()

	err := godotenv.Load("/home/ashith/Hacfy/it_management_system/.env")
	if err != nil {
		log.Fatal("error while loading env %s", err)
	}

	db := NewDatabase()
	db.CheckStatus()
	// defer db.Close()

	Start(db.db, addr)
}
