package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserTokenModel struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserType  string `json:"user_type"`
	jwt.RegisteredClaims
}

type UserInterface interface {
	UserLogin(echo.Context) (int, string, string, string, error)
}
