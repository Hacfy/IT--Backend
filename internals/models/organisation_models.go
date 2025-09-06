package models

import "github.com/labstack/echo/v4"

type OrganisationModel struct {
	OrganisationID          int    `json:"organisation_id"`
	OrganisationName        string `json:"organisation_name"`
	OrganisationEmail       string `json:"organisation_email"`
	OrganisationPhoneNumber string `json:"organisation_phone_number"`
	OrganisationPassword    string `json:"organisation_password"`
	OrganisationMainAdminID int    `json:"organisation_main_admin_id"`
}

type CreateOrganisationModel struct {
	OrganisationName           string `json:"organisation_name" validate:"required"`
	OrganisationEmail          string `json:"organisation_email" validate:"required,email"`
	OrganisationPhoneNumber    string `json:"organisation_phone_number" validate:"required"`
	CreateOrganisationPassword string `json:"create_organisation_password" validate:"required"`
}

type DeleteOrganisationModel struct {
	OrganisationEmail          string `json:"organisation_email" validate:"required,email"`
	DeleteOrganisationPassword string `json:"delete_organisation_password" validate:"required"`
	OrganisationID             int    `json:"organisation_id" validate:"required"`
}

type GetAllOrganisationsModel struct {
	OrganisationID          int    `json:"organisation_id"`
	OrganisationName        string `json:"organisation_name"`
	OrganisationEmail       string `json:"organisation_email"`
	OrganisationPhoneNumber string `json:"organisation_phone_number"`
}

type OrganisationInterface interface {
	CreateSuperAdmin(echo.Context) (int, error)
	DeleteSuperAdmin(echo.Context) (int, error)
	GetAllSuperAdmins(echo.Context) (int, []AllSuperAdminsDetailsModel, error)
	ReassignSuperAdmin(echo.Context) (int, error)
}

// type GetAllOrgDepartmentsModel struct {
// 	BranchID           int    `json:"branch_id" validate:"required"`
// 	BranchName         string `json:"branch_name" validate:"required"`
// 	OrganisationID     int    `json:"organisation_id" validate:"required"`
// 	OrganisationName   int    `json:"organisation_name" validate:"required"`
// 	DepartmentName     string `json:"department_name" validate:"required"`
// 	DepartmentID       int    `json:"department_id" validate:"required"`
// 	DepartmentHeadName string `json:"department_head_name" validate:"required"`
// 	NoOfWorkspaces     int    `json:"no_of_workspaces" validate:"required"`
// 	Issues             int    `json:"issues" validate:"required"`
// }
