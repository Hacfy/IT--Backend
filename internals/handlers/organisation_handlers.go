package handlers

import (
	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type OrganisationHandler struct {
	OrgRepo models.OrganisationInterface
}

func (oh *OrganisationHandler) CreateSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.CreateSuperAdmin(e)
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
