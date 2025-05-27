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
	authHandler := handlers.NewAuthHandler(repository.NewAuthRepo(db))
	organisatoinHandler := handlers.NewOrganisationHandler(repository.NewOrgRepo(db))
	superAdminHandler := handlers.NewSuperAdminHandler(repository.NewSuperAdminRepo(db))

	e.POST("/main_admin/register", mainAdminHandler.CreateMainAdminHandler)
	e.POST("/main_admin/login", mainAdminHandler.LoginMainAdminHandler)
	e.POST("/main_admin/create/organisation", mainAdminHandler.CreateOrganisationHandler)

	e.POST("/auth/login/users", authHandler.UserLoginHandler)

	e.POST("/organisation/create/superAdmin", organisatoinHandler.CreateSuperAdminHandler)
	e.DELETE("/organisation/delete/superAdmin", organisatoinHandler.DeleteSuperAdminHandler)

	e.POST("/superAdmin/create/branch", superAdminHandler.CreateBranchHandler)
	return e
}
