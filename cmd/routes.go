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

	e.Use(middleware.Logger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	mainAdminHandler := handlers.NewMainAdmin_Handler(repository.NewMainAdminRepo(db))

	e.POST("/main_admin/register", mainAdminHandler.CreateMainAdminHandler)
	e.POST("/main_admin/login", mainAdminHandler.LoginMainAdminHandler)
	e.DELETE("/main_admin/delete/main_admin", mainAdminHandler.DeleteMainAdminHandler)
	e.POST("/main_admin/create/organisation", mainAdminHandler.CreateOrganisationHandler)
	e.DELETE("/main_admin/delete/organisation", mainAdminHandler.DeleteOrganisationHandler)

	authHandler := handlers.NewAuthHandler(repository.NewAuthRepo(db))

	e.POST("/auth/login/users", authHandler.UserLoginHandler)
	e.POST("/auth/logout/users", authHandler.UserLogoutHandler)
	e.POST("/auth/change/password", authHandler.ChangeUserPasswordHandler)

	organisatoinHandler := handlers.NewOrganisationHandler(repository.NewOrgRepo(db))

	e.POST("/organisation/create/superAdmin", organisatoinHandler.CreateSuperAdminHandler)
	e.DELETE("/organisation/delete/superAdmin", organisatoinHandler.DeleteSuperAdminHandler)

	superAdminHandler := handlers.NewSuperAdminHandler(repository.NewSuperAdminRepo(db))

	e.POST("/superAdmin/create/branch", superAdminHandler.CreateBranchHandler)
	e.POST("/superAdmin/update/branchHead", superAdminHandler.UpdateBranchHeadHandler)
	e.DELETE("/superAdmin/delete/branch", superAdminHandler.DeleteBranchHandler)

	branchHandler := handlers.NewBranchHandler(repository.NewBranchRepo(db))

	e.POST("/branch/create/department", branchHandler.CreateDepartmentHandler)
	e.POST("/branch/update/departmentHead", branchHandler.UpdateDepartmentHeadHandler)
	e.POST("/branch/create/warehouse", branchHandler.CreateWarehouseHandler)
	e.POST("/branch/update/warehouseHead", branchHandler.UpdateWarehouseHeadHandler)

	departmentHandler := handlers.NewDepartmentHandler(repository.NewDepartmentRepo(db))

	e.POST("/department/create/workspace", departmentHandler.CreateWorkspaceHandler)
	e.DELETE("/department/delete/workspace", departmentHandler.DeleteWorkspaceHandler)
	e.POST("/department/raise/issue", departmentHandler.RaiseIssueHandler)

	warehouseHandler := handlers.NewWarehouse_Handler(repository.NewWarehouseRepo(db))

	e.POST("/warehouse/create/component", warehouseHandler.CreateComponentHandler)
	e.DELETE("/warehouse/delete/component", warehouseHandler.DeleteComponentHandler)
	e.POST("/warehouse/add/component/units", warehouseHandler.AddComponentUnitsHandler)
	e.POST("/warehouse/assign/units", warehouseHandler.AssignUnitsHandler)
	e.GET("/warehouse/get/all/issues", warehouseHandler.GetAllIssuesHandler)

	return e
}
