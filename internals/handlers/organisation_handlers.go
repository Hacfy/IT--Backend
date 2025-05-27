package handlers

import (
	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type OrganisationHandler struct {
	OrgRepo models.OrganisationInterface
}

func NewOrganisationHandler(organisaion models.OrganisationInterface) *OrganisationHandler {
	return &OrganisationHandler{
		OrgRepo: organisaion,
	}
}

func (oh *OrganisationHandler) CreateSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.CreateSuperAdmin(e)
	if err != nil {
		return e.JSON(status, echo.Map{
			"error": err.Error(),
		})
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (oh *OrganisationHandler) DeleteSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.DeleteSuperAdmin(e)
	if err != nil {
		return e.JSON(status, echo.Map{
			"error": err.Error(),
		})
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}
