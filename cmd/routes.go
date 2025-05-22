package main

import (
	"database/sql"

	"github.com/Hacfy/IT_INVENTORY/internals/handlers"
	"github.com/Hacfy/IT_INVENTORY/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitialiseHttpRouter(db *sql.DB) *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	mainAdminHandler := handlers.NewMainAdmin_Handler(repository.NewMainAdminRepo(db))
	e.POST("/main_admin/register", mainAdminHandler.CreateMainAdminHandler)
	e.POST("/main_admin/login", mainAdminHandler.LoginMainAdminHandler)
	e.POST("/main_admin/create/organisation", mainAdminHandler.CreateOrganisationHandler)
	return e
}
