package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserLoginModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLogoutModel struct {
	Email string `json:"email" validate:"required,email"`
}

type UserTokenModel struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
	UserType  string `json:"user_type"`
	jwt.RegisteredClaims
}

type UserInterface interface {
	UserLogin(echo.Context) (int, string, string, string, error)
	ChangeUserPassword(e echo.Context) (int, string, string, string, error)
	UserLogout(e echo.Context) (int, error)
	ResetPassword(e echo.Context) (int, string, string, string, error)
	ForgotPassword(e echo.Context) (int, time.Time, error)
	VerifyForgotPasswordRequest(e echo.Context) (int, error)
}

type ChangePasswordModel struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type CookieModel struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"user_email" validate:"required,email"`
	UserType  string `json:"user_type"`
	jwt.RegisteredClaims
}

type ForgotPasswordModel struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyForgotPasswordRequestModel struct {
	Email string    `json:"email" validate:"required,email"`
	Otp   string    `json:"otp" validate:"required"`
	Time  time.Time `json:"time" validate:"required"`
}

type ResetPasswordModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
