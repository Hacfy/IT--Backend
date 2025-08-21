package handlers

import (
	"net/http"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type MainAdminHandler struct {
	MainAdminRepo models.MainAdminInterface
}

func NewMainAdmin_Handler(mainAdmin models.MainAdminInterface) *MainAdminHandler {
	return &MainAdminHandler{
		MainAdminRepo: mainAdmin,
	}
}

func (ma *MainAdminHandler) CreateMainAdminHandler(e echo.Context) error {
	status, err := ma.MainAdminRepo.CreateMainAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (ma *MainAdminHandler) LoginMainAdminHandler(e echo.Context) error {
	status, accessToken, refreshToken, token, err := ma.MainAdminRepo.LoginMainAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = accessToken
	accessCookie.HttpOnly = true
	accessCookie.Secure = false
	accessCookie.Expires = time.Now().Add(15 * time.Hour)
	e.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = refreshToken
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = false
	refreshCookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	e.SetCookie(refreshCookie)

	return e.JSON(status, echo.Map{
		"message": "successfull",
		"token":   token,
	})
}

func (ma *MainAdminHandler) CreateOrganisationHandler(e echo.Context) error {
	status, err := ma.MainAdminRepo.CreateOrganisation(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (ma *MainAdminHandler) DeleteMainAdminHandler(e echo.Context) error {
	status, err := ma.MainAdminRepo.DeleteMainAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (ma *MainAdminHandler) DeleteOrganisationHandler(e echo.Context) error {
	status, err := ma.MainAdminRepo.DeleteOrganisation(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (ma *MainAdminHandler) GetAllOrganisationsHandler(e echo.Context) error {
	status, orgs, err := ma.MainAdminRepo.GetAllOrganisations(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message":             "successfull",
		"organisations":       orgs,
		"no_of_organisations": len(orgs),
	})
}

func (ma *MainAdminHandler) GetAllMainAdminsHandler(e echo.Context) error {
	status, main_admins, err := ma.MainAdminRepo.GetAllMainAdmins(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message":           "successfull",
		"main_admins":       main_admins,
		"no_of_main_admins": len(main_admins),
	})
}
