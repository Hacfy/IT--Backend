package handlers

import (
	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	OrgRepo models.OrganizationInterface
}

func NeworganizationHandler(organisaion models.OrganizationInterface) *OrganizationHandler {
	return &OrganizationHandler{
		OrgRepo: organisaion,
	}
}

func (oh *OrganizationHandler) CreateSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.CreateSuperAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (oh *OrganizationHandler) DeleteSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.DeleteSuperAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (oh *OrganizationHandler) GetAllSuperAdminsHandler(e echo.Context) error {
	status, superAdmins, err := oh.OrgRepo.GetAllSuperAdmins(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"superAdmins": superAdmins,
		"total":       len(superAdmins),
	})
}

func (oh *OrganizationHandler) ReassignSuperAdminHandler(e echo.Context) error {
	status, err := oh.OrgRepo.ReassignSuperAdmin(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}
