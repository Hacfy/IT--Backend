package models

import (
	"github.com/labstack/echo/v4"
)

type WarehouseModel struct {
	WarehouseID        int    `json:"warehouse_id"`
	WarehouseUserName  string `json:"warehouse_user_name"`
	WarehouseUserEmail string `json:"warehouse_user_email"`
	BranchID           int    `json:"branch_id"`
}

type CreateWarehouseModel struct {
	WarehouseUserName  string `json:"warehouse_user_name" validate:"required"`
	WarehouseUserEmail string `json:"warehouse_user_email" validate:"required,email"`
}

type UpdateWarehouseHeadModel struct {
	WarehouseID           int    `json:"warehouse_id" validate:"required"`
	WarehouseHeadEmail    string `json:"warehouse_head_email" validate:"required,email"`
	NewWarehouseHeadName  string `json:"new_warehouse_head_name" validate:"required"`
	NewWarehouseHeadEmail string `json:"new_warehouse_head_email" validate:"required,email"`
}

type WarehouseInterface interface {
	CreateComponent(echo.Context) (int, string, error)
	DeleteComponent(echo.Context) (int, error)
	AddComponentUnits(echo.Context) (int, error)
	AssignUnits(echo.Context) (int, error)
	GetAllIssues(echo.Context) (int, []IssueModel, int, int, int, error)
}

type GetAllIssuesModel struct {
	IssueID      int    `json:"issue_id" validate:"required"`
	Issue        string `json:"issue" validate:"required"`
	DepartmentID int    `json:"department_id" validate:"required"`
	WorkspaceID  int    `json:"workspace_id" validate:"required"`
}
