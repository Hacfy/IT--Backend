package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hacfy/IT_INVENTORY/pkg/database"
)

func Start(db *sql.DB, addr *string) {
	e := InitialiseHttpRouter(db)

	query := database.NewDBinstance(db)

	err := query.InitialiseDBqueries()
	if err != nil {
		log.Fatalf("Unable to Initialize Database %v", err)
	}

	log.Printf("Starting server at %s", *addr)
	go func() {
		e.Start(*addr)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		db.Close()
		log.Println("db closed")
		cancel()
		log.Println("ctx cancled")
	}()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	fmt.Println("Server exited properly.")
}
