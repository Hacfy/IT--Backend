package models

import "github.com/labstack/echo/v4"

type BranchModel struct {
	BranchID       int    `json:"branch_id"`
	OrgID          int    `json:"org_id"`
	SuperAdminID   int    `json:"super_admin_id"`
	BranchName     string `json:"branch_name"`
	BranchLocation string `json:"branch_location"`
}

type BranchHeadModel struct {
	BranchID           int    `json:"branch_id"`
	BranchHeadID       int    `json:"branch_head_id"`
	BranchHeadName     string `json:"branch_head_name"`
	BranchHeadEmail    string `json:"branch_head_email"`
	BranchHeadPassword string `json:"branch_head_password"`
}

type CreateBranchModel struct {
	BranchName      string `json:"branch_name" validate:"required"`
	BranchLocation  string `json:"branch_location" validate:"required"`
	BranchHeadName  string `json:"branch_head_name" validate:"required"`
	BranchHeadEmail string `json:"branch_head_email" validate:"required,email"`
}

type DeleteBranchModel struct {
	BranchID  int    `json:"branch_id" validate:"required"`
	BrachName string `json:"branch_name" validate:"required"`
}

type UpdateBranchHeadModel struct {
	BranchHeadID       int    `json:"branch_head_id" validate:"required"`
	BranchHeadEmail    string `json:"branch_head_email" validate:"required,email"`
	NewBranchHeadName  string `json:"new_branch_head_name" validate:"required"`
	NewBranchHeadEmail string `json:"new_branch_head_email" validate:"required,email"`
}

type BranchInterface interface {
	CreateDepartment(echo.Context) (int, error)
	UpdateDepartmentHead(echo.Context) (int, error)
	CreateWarehouse(echo.Context) (int, error)
	UpdateWarehouseHead(echo.Context) (int, error)
}

type GetAllBranchesModel struct {
	BranchName     string `json:"branch_name" validate:"required"`
	BranchID       int    `json:"branch_id" validate:"required"`
	BranchLocation string `json:"branch_location" validate:"required"`
}
