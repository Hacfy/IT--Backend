package models

import "github.com/labstack/echo/v4"

type DepartmentModel struct {
	DepartmentID int    `json:"department_id"`
	BranchID     int    `json:"branch_id"`
	BranchName   string `json:"branch_name"`
}

type CreateDepartmentModel struct {
	DepartmentName      string `json:"department_id" validate:"required"`
	DepartmentHeadName  string `json:"department_head_name" validate:"required"`
	DepartmentHeadEmail string `json:"department_head_email" validate:"required,email"`
}

type UpdateDepartmentHeadModel struct {
	DepartmentID           int    `json:"department_id" validate:"required"`
	DepartmentHeadEmail    string `json:"department_head_email" validate:"required,email"`
	NewDepartmentHeadName  string `json:"new_department_head_name" validate:"required"`
	NewDepartmentHeadEmail string `json:"new_department_head_email" validate:"required,email"`
}

type DepartmentInterface interface {
	CreateWorkspace(echo.Context) (int, int, error)
	DeleteWorkspace(echo.Context) (int, error)
	RaiseIssue(echo.Context) (int, int, error)
}

type GetAllDepartmentsModel struct {
	DepartmentName     string `json:"department_name" validate:"required"`
	DepartmentID       int    `json:"department_id" validate:"required"`
	DepartmentHeadName string `json:"department_head_name" validate:"required"`
	NoOfWorkspaces     int    `json:"no_of_workspaces" validate:"required"`
	Issues             int    `json:"issues" validate:"required"`
}

type RaiseIssueModel struct {
	DepartmentID int    `json:"department_id" validate:"required"`
	WarehouseID  int    `json:"warehouse_id" validate:"required"`
	WorkspaceID  int    `json:"workspace_id" validate:"required"`
	UnitID       int    `json:"unit_id" validate:"required"`
	UnitPrefix   string `json:"unit_prefix" validate:"required"`
	Issue        string `json:"issue" validate:"required"`
}
