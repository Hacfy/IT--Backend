package middleware

import (
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := c.Request().Header.Get("Authorization")

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			return c.JSON(401, echo.Map{"error": "Unauthorized"})
		}

		c.Set("userID", claims.UserID)
		c.Set("userType", claims.UserType)
		c.Set("userEmail", claims.UserEmail)
		return next(c)
	}
}

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	roleMap := make(map[string]bool)
	for _, r := range allowedRoles {
		roleMap[r] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("userType").(string)
			if !ok || !roleMap[role] {
				return c.JSON(403, echo.Map{"error": "Access denied"})
			}
			return next(c)
		}
	}
}
