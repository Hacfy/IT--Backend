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

type GetAllSuperAdminsModel struct {
	OrganizationID int `json:"organization_id" validate:"required"`
}

type AllSuperAdminsDetailsModel struct {
	SuperAdminID    int    `json:"super_admin_id"`
	SuperAdminName  string `json:"super_admin_name"`
	SuperAdminEmail string `json:"super_admin_email"`
	NoOfBranches    int    `json:"no_of_branches"`
}

type GetSuperAdminDetailsModel struct {
	SuperAdminID int `json:"super_admin_id" validate:"required"`
}

type SuperAdminDetailsModel struct {
	SuperAdminID           int    `json:"super_admin_id"`
	SuperAdminName         string `json:"super_admin_name"`
	SuperAdminEmail        string `json:"super_admin_email"`
	NoOfBranches           int    `json:"no_of_branches"`
	NoOfDepartments        int    `json:"no_of_departments"`
	NoOfWarehouses         int    `json:"no_of_warehouses"`
	NoOfWorkspaces         int    `json:"no_of_workspaces"`
	TotalTypesOfComponents int    `json:"total_types_of_components"`
	TotalUnitsOfComponents int    `json:"total_units_of_components"`
	TotalNoOfIssues        int    `json:"total_no_of_issues"`
}

type ReassignSuperAdminModel struct {
	OldSuperAdminID int `json:"old_super_admin_id" validate:"required"`
	NewSuperAdminID int `json:"new_super_admin_id" validate:"required"`
}
