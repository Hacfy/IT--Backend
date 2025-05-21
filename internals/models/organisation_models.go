package models

import "github.com/labstack/echo/v4"

type OrganisationModel struct {
	OrganisationID          int    `json:"organisation_id"`
	OrganisationName        string `json:"organisation_name"`
	OrganisationEmail       string `json:"organisation_email"`
	OrganisationPassword    string `json:"organisation_password"`
	OrganisationMainAdminID int    `json:"organisation_main_admin_id"`
}

type CreateOrganisationModel struct {
	OrganisationName  string `json:"organisation_name" validate:"required"`
	OrganisationEmail string `json:"organisation_email" validate:"required,email"`
}

type OrganisationInterface interface {
	CreateSuperAdmin(echo.Context) (int, error)
}
