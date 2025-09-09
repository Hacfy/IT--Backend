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
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	mainAdminHandler := handlers.NewMainAdmin_Handler(repository.NewMainAdminRepo(db))

	e.POST("/main_admin/register", mainAdminHandler.CreateMainAdminHandler)                 //
	e.POST("/main_admin/login", mainAdminHandler.LoginMainAdminHandler)                     //
	e.DELETE("/main_admin/delete/main_admin", mainAdminHandler.DeleteMainAdminHandler)      //
	e.POST("/main_admin/create/organisation", mainAdminHandler.CreateOrganisationHandler)   //
	e.DELETE("/main_admin/delete/organisation", mainAdminHandler.DeleteOrganisationHandler) //
	e.GET("/main_admin/get/all/organisations", mainAdminHandler.GetAllOrganisationsHandler) //
	e.GET("/main_admin/get/all/main_admins", mainAdminHandler.GetAllMainAdminsHandler)      //

	authHandler := handlers.NewAuthHandler(repository.NewAuthRepo(db))

	e.POST("/auth/login/users", authHandler.UserLoginHandler)                              //
	e.POST("/auth/logout/users", authHandler.UserLogoutHandler)                            //
	e.PUT("/auth/change/password", authHandler.ChangeUserPasswordHandler)                  //
	e.POST("/auth/reset/password", authHandler.ResetPasswordHandler)                       //
	e.POST("/auth/forgot/password", authHandler.ForgotPasswordHandler)                     //
	e.POST("/auth/verify/forgot/password", authHandler.VerifyForgotPasswordRequestHandler) //

	organisatoinHandler := handlers.NewOrganisationHandler(repository.NewOrgRepo(db))

	e.POST("/organisation/create/superAdmin", organisatoinHandler.CreateSuperAdminHandler)              //
	e.DELETE("/organisation/delete/superAdmin", organisatoinHandler.DeleteSuperAdminHandler)            //
	e.GET("/organisation/get/all/superAdmins", organisatoinHandler.GetAllSuperAdminsHandler)            //
	e.POST("/organisation/reassign/superAdmin/branches", organisatoinHandler.ReassignSuperAdminHandler) //
	//GET /organisation/get/details

	superAdminHandler := handlers.NewSuperAdminHandler(repository.NewSuperAdminRepo(db))

	e.POST("/superAdmin/create/branch", superAdminHandler.CreateBranchHandler)        //
	e.PUT("/superAdmin/update/branchHead", superAdminHandler.UpdateBranchHeadHandler) //
	e.DELETE("/superAdmin/delete/branch", superAdminHandler.DeleteBranchHandler)      //
	// GET /superAdmin/get/branch/details

	branchHandler := handlers.NewBranchHandler(repository.NewBranchRepo(db))

	e.POST("/branch/create/department", branchHandler.CreateDepartmentHandler)        //
	e.PUT("/branch/update/departmentHead", branchHandler.UpdateDepartmentHeadHandler) //
	e.POST("/branch/create/warehouse", branchHandler.CreateWarehouseHandler)          //
	e.PUT("/branch/update/warehouseHead", branchHandler.UpdateWarehouseHeadHandler)   //
	e.DELETE("/branch/delete/department", branchHandler.DeleteDepartmentHandler)      //
	e.DELETE("/branch/delete/warehouse", branchHandler.DeleteWarehouseHandler)        //
	// GET /branch/get/department/details
	// GET /branch/get/warehouse/details

	departmentHandler := handlers.NewDepartmentHandler(repository.NewDepartmentRepo(db))

	e.POST("/department/create/workspace", departmentHandler.CreateWorkspaceHandler)   //
	e.DELETE("/department/delete/workspace", departmentHandler.DeleteWorkspaceHandler) //
	e.POST("/department/raise/issue", departmentHandler.RaiseIssueHandler)             //
	e.POST("/department/request/new/units", departmentHandler.RequestNewUnitsHandler)  //
	e.GET("/department/get/all/requests", departmentHandler.GetAllDepartmentRequestsHandler)
	e.GET("/department/get/request/details", departmentHandler.GetDepartmentRequestDetailsHandler)

	warehouseHandler := handlers.NewWarehouse_Handler(repository.NewWarehouseRepo(db))

	e.POST("/warehouse/create/component", warehouseHandler.CreateComponentHandler)
	e.DELETE("/warehouse/delete/component", warehouseHandler.DeleteComponentHandler)
	e.POST("/warehouse/add/component/units", warehouseHandler.AddComponentUnitsHandler)
	e.PATCH("/warehouse/assign/units", warehouseHandler.AssignUnitsHandler)
	e.GET("/warehouse/get/all/issues", warehouseHandler.GetAllIssuesHandler)
	e.GET("/warehouse/get/all/components", warehouseHandler.GetAllWarehouseComponentsHandler)
	e.GET("/warehouse/get/all/component/units", warehouseHandler.GetAllWarehouseComponentUnitsHandler)
	e.GET("/warehouse/get/issue/details", warehouseHandler.GetIssueDetailsHandler)
	e.GET("/warehouse/get/unit/history", warehouseHandler.GetUnitAssignmentHistoryHandler)
	e.PUT("/warehouse/update/issue/status", warehouseHandler.UpdateIssueStatusHandler)
	e.PUT("/warehouse/update/component/name", warehouseHandler.UpdateComponentNameHandler)
	e.GET("/warehouse/get/assigned/units", warehouseHandler.GetAssignedUnitsHandler)
	e.PUT("/warehouse/update/component/unit/maintainance", warehouseHandler.UpdateMaintenanceCostHandler)
	// GET /warehouse/get/component/details

	detailsHandler := handlers.NewDetailsHandler(repository.NewDetailsRepo(db))

	e.GET("/details/get/all/departments", detailsHandler.GetAllDepartmentsHandler)
	e.GET("/details/get/all/departments/issues", detailsHandler.GetDepartmentIssuesHandler)
	e.GET("/details/get/all/departments/workspaces", detailsHandler.GetDepartmentWorkspacesHandler)
	e.GET("/details/get/all/branches", detailsHandler.GetAllBranchesHandler)
	e.GET("/details/get/all/warehouses", detailsHandler.GetAllWarehousesHandler)
	e.GET("/details/get/all/department/outOfWarentyUnits", detailsHandler.GetAllDepartmentOutOfWarentyUnitsHandler)
	e.GET("/details/get/all/warehouse/outOfWarentyUnits", detailsHandler.GetAllOutOfWarentyUnitsInWarehouseHandler)

	excelHandler := handlers.NewExcelHandler(repository.NewExcelRepo(db))

	e.GET("/excel/download/component/maintainance/report", excelHandler.DownloadComponentMaintainanceReportHandler)

	return e
}
