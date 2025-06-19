package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type MainAdminModel struct {
	MainAdminID       int    `json:"main_admin_id"`
	MainAdminEmail    string `json:"main_admin_email"`
	MainAdminPassword string `json:"main_admin_password"`
}

type CreateMainAdminModel struct {
	MainAdminEmail  string `json:"main_admin_email"  validate:"required,email"`
	CompanyPassword string `json:"company_password" validate:"required"`
}

type LoginMainAdminModel struct {
	MainAdminEmail    string `json:"main_admin_email" validate:"required,email"`
	MainAdminPassword string `json:"main_admin_password" validate:"required"`
}

type DeleteMainAdminModel struct {
	MainAdminEmail  string `json:"main_admin_email" validate:"required,email"`
	CompanyPassword string `json:"company_password" validate:"required"`
}

type GetAllMainAdminModel struct {
	CompanyPassword string `json:"company_password" validate:"required"`
}

type MainAdminInterface interface {
	CreateMainAdmin(echo.Context) (int, error)
	CreateOrganisation(echo.Context) (int, error)
	LoginMainAdmin(echo.Context) (int, string, string, string, error)
	DeleteMainAdmin(echo.Context) (int, error)
	DeleteOrganisation(echo.Context) (int, error)
}

type MainAdminTokenModel struct {
	MainAdminID    int    `json:"main_admin_id"`
	MainAdminEmail string `json:"main_admin_email"`
	jwt.RegisteredClaims
}
