package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	addr := flag.String("addr", ":8080", "Address at which the server will run")
	flag.Parse()

	err := godotenv.Load("/home/ashith/Hacfy/IT_INVENTORY/.env.local")
	if err != nil {
		log.Fatalf("error while loading env %s", err)
	}

	db := NewDatabase()
	log.Println("connected to db")
	db.CheckStatus()
	defer db.Close()

	Start(db.db, addr)
}
