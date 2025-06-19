package models

import "github.com/labstack/echo/v4"

type SuperAdminModel struct {
	SuperAdminID       int    `json:"super_admin_id"`
	Org_ID             int    `json:"org_id"`
	SuperAdminName     string `json:"super_admin_name"`
	SuperAdminEmail    string `json:"super_admin_email"`
	SuperAdminPassword string `json:"super_admin_password"`
}

type CreateSuperAdminModel struct {
	SuperAdminName  string `json:"super_admin_name" validate:"required"`
	SuperAdminEmail string `json:"super_admin_email" validate:"required,email"`
}

//login in auth model

type DeleteSuperAdminModel struct {
	SuperAdminEmail string `json:"super_admin_email" validate:"required,email"`
}

type SuperAdminInterface interface {
	CreateBranch(echo.Context) (int, error)
	DeleteBranch(echo.Context) (int, error)
	UpdateBranchHead(echo.Context) (int, error)
}
