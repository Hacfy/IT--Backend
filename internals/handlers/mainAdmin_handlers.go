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
		return e.JSON(status, echo.Map{
			"error":   err.Error(),
			"message": "unsuccessfull",
		})
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (ma *MainAdminHandler) LoginMainAdminHandler(e echo.Context) error {
	status, accessToken, refreshToken, token, err := ma.MainAdminRepo.LoginMainAdmin(e)
	if err != nil {
		return e.JSON(status, echo.Map{
			"error":   err.Error(),
			"message": "unsuccessfull",
		})
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
	refreshCookie.Expires = time.Now().Add(5 * 24 * time.Hour)
	e.SetCookie(refreshCookie)

	return e.JSON(status, echo.Map{
		"message": "successfull",
		"token":   token,
	})
}

func (ma *MainAdminHandler) CreateOrganisationHandler(e echo.Context) error {
	status, err := ma.MainAdminRepo.CreateOrganisation(e)
	if err != nil {
		return e.JSON(status, echo.Map{
			"message": "unsuccessfull",
			"error":   err.Error(),
		})
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}
