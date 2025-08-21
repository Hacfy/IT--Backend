package models

import (
	"github.com/labstack/echo/v4"
)

type AllDepartmentsModel struct {
	DepartmentName     string `json:"department_name" validate:"required"`
	DepartmentID       int    `json:"department_id" validate:"required"`
	DepartmentHeadName string `json:"department_head_name" validate:"required"`
	NoOfWorkspaces     int    `json:"no_of_workspaces" validate:"required"`
	Issues             int    `json:"issues" validate:"required"`
}

type GetAllDepartmentsModel struct {
	BranchID int `json:"branch_id" validate:"required"`
}

type DetailsInterface interface {
	GetAllDepartmentsRepo(echo.Context) ([]AllDepartmentsModel, int, int, int, int, error)
	GetDepartmentIssues(echo.Context) (int, []DepartmentIssuesModel, int, int, int, error)
	GetDepartmentWorkspaces(echo.Context) ([]DepartmentWorkspaceModel, int, int, int, int, error)
	GetAllBranches(echo.Context) ([]AllBranchesModel, int, int, int, int, error)
}

type DepartmentWorkspaceModel struct {
	WorkspaceID   int    `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}

type AllBranchesModel struct {
	BranchID       int    `json:"branch_id"`
	BranchName     string `json:"branch_name"`
	BranchLocation string `json:"branch_location"`
	BranchHeadName string `json:"branch_head_name"`
}

type BranchDetailsModel struct {
	BranchID               int    `json:"branch_id"`
	BranchName             string `json:"branch_name"`
	BranchLocation         string `json:"branch_location"`
	BranchHeadName         string `json:"branch_head_name"`
	BranchHeadEmail        string `json:"branch_head_email"`
	NoOfDepartments        int    `json:"no_of_departments"`
	NoOfWorkspaces         int    `json:"no_of_workspaces"`
	TotalTypesOfComponents int    `json:"total_types_of_components"`
	TotalUnitsOfComponents int    `json:"total_units_of_components"`
	TotalNoOfIssues        int    `json:"total_no_of_issues"`
}

type GetAllWarehousesModel struct {
	BranchID int `json:"branch_id" validate:"required"`
}

type AllWarehousesModel struct {
	WarehouseID            int    `json:"warehouse_id"`
	Warehouse              string `json:"warehouse"`
	TotalTypesOfComponents int    `json:"total_types_of_components"`
	TotalUnitsOfComponents int    `json:"total_units_of_components"`
	TotalNoOfIssues        int    `json:"total_no_of_issues"`
}
