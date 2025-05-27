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
}

type OrganisationInterface interface {
	CreateSuperAdmin(echo.Context) (int, error)
	DeleteSuperAdmin(echo.Context) (int, error)
}
