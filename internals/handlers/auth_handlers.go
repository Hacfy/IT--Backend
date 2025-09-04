package handlers

import (
	"net/http"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthRepo models.UserInterface
}

func NewAuthHandler(user models.UserInterface) *AuthHandler {
	return &AuthHandler{
		AuthRepo: user,
	}
}

func (ah *AuthHandler) UserLoginHandler(e echo.Context) error {
	status, accessToken, refreshToken, token, err := ah.AuthRepo.UserLogin(e)
	if err != nil {
		if status == http.StatusFound {
			return e.Redirect(http.StatusFound, "/change-password")
		}
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
		"message": "logged in successfully",
		"token":   token,
	})
}

func (ah *AuthHandler) ChangeUserPasswordHandler(e echo.Context) error {
	status, accessToken, refreshToken, token, err := ah.AuthRepo.ChangeUserPassword(e)
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

func (ah *AuthHandler) UserLogoutHandler(e echo.Context) error {
	status, err := ah.AuthRepo.UserLogout(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = ""
	accessCookie.HttpOnly = true
	accessCookie.Secure = false
	accessCookie.Expires = time.Time{}
	e.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = ""
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = false
	refreshCookie.Expires = time.Time{}
	e.SetCookie(refreshCookie)

	return e.JSON(status, echo.Map{
		"message": "logged out successfully",
	})
}

func (ah *AuthHandler) ResetPasswordHandler(e echo.Context) error {
	status, accessToken, refreshToken, token, err := ah.AuthRepo.ResetPassword(e)
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

func (ah *AuthHandler) ForgotPasswordHandler(e echo.Context) error {
	status, Time, err := ah.AuthRepo.ForgotPassword(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
		"time":    Time,
	})
}

func (ah *AuthHandler) VerifyForgotPasswordRequestHandler(e echo.Context) error {
	status, err := ah.AuthRepo.VerifyForgotPasswordRequest(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}
