package models

import "github.com/labstack/echo/v4"

type OrganizationModel struct {
	OrganizationID          int    `json:"organization_id"`
	OrganizationName        string `json:"organization_name"`
	OrganizationEmail       string `json:"organization_email"`
	OrganizationPhoneNumber string `json:"organization_phone_number"`
	OrganizationPassword    string `json:"organization_password"`
	OrganizationMainAdminID int    `json:"organization_main_admin_id"`
}

type CreateOrganizationModel struct {
	OrganizationName           string `json:"organization_name" validate:"required"`
	OrganizationEmail          string `json:"organization_email" validate:"required,email"`
	OrganizationPhoneNumber    string `json:"organization_phone_number" validate:"required"`
	CreateOrganizationPassword string `json:"create_Organization_password" validate:"required"`
}

type DeleteOrganizationModel struct {
	OrganizationEmail          string `json:"organization_email" validate:"required,email"`
	DeleteOrganizationPassword string `json:"delete_Organization_password" validate:"required"`
	OrganizationID             int    `json:"organization_id" validate:"required"`
}

type GetAllOrganizationModel struct {
	OrganizationID          int    `json:"organization_id"`
	OrganizationName        string `json:"organization_name"`
	OrganizationEmail       string `json:"organization_email"`
	OrganizationPhoneNumber string `json:"organization_phone_number"`
}

type OrganizationInterface interface {
	CreateSuperAdmin(echo.Context) (int, error)
	DeleteSuperAdmin(echo.Context) (int, error)
	GetAllSuperAdmins(echo.Context) (int, []AllSuperAdminsDetailsModel, error)
	ReassignSuperAdmin(echo.Context) (int, error)
}

// type GetAllOrgDepartmentsModel struct {
// 	BranchID           int    `json:"branch_id" validate:"required"`
// 	BranchName         string `json:"branch_name" validate:"required"`
// 	OrganizationID     int    `json:"Organization_id" validate:"required"`
// 	OrganizationName   int    `json:"Organization_name" validate:"required"`
// 	DepartmentName     string `json:"department_name" validate:"required"`
// 	DepartmentID       int    `json:"department_id" validate:"required"`
// 	DepartmentHeadName string `json:"department_head_name" validate:"required"`
// 	NoOfWorkspaces     int    `json:"no_of_workspaces" validate:"required"`
// 	Issues             int    `json:"issues" validate:"required"`
// }
