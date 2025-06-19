package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/labstack/echo/v4"
)

func Start(db *sql.DB, addr *string) {
	e := InitialiseHttpRouter(db)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code                = http.StatusInternalServerError
			message interface{} = "Something went wrong"
		)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message
		}

		c.JSON(code, echo.Map{
			"error": message,
		})
	}

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
