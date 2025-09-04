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
		return echo.NewHTTPError(status, err.Error())
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (oh *OrganisationHandler) DeleteSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.DeleteSuperAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (oh *OrganisationHandler) GetAllSuperAdminsHandler(e echo.Context) error {
	status, superAdmins, err := oh.OrgRepo.GetAllSuperAdmins(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"superAdmins": superAdmins,
		"total":       len(superAdmins),
	})
}

func (oh *OrganisationHandler) ReassignSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.ReassignSuperAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}
