package middleware

import (
	"net/http"
	"os"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CookieMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			accessClaims := &models.CookieModel{}
			accessCookie, err := c.Request().Cookie("access_token")

			if err == nil {
				token, err := parseJWTFromCookie(accessCookie.Value, accessClaims)
				if err == nil && token.Valid && accessClaims.Issuer == "IT_INVENTORY" {
					return next(c)
				}
			}

			refreshClaims := &models.CookieModel{}
			refreshCookie, err := c.Request().Cookie("refresh_token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: No refresh cookie")
			}

			token, err := parseJWTFromCookie(refreshCookie.Value, refreshClaims)
			if err != nil || !token.Valid || refreshClaims.Issuer != "IT_INVENTORY" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
			}

			return next(c)
		}
	}
}

func parseJWTFromCookie(tokenStr string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
