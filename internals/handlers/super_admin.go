package handlers

import (
	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type SuperAdminHandler struct {
	SuperAdminRepo models.SuperAdminInterface
}

func NewSuperAdminHandler(superAdmin models.SuperAdminInterface) *SuperAdminHandler {
	return &SuperAdminHandler{
		SuperAdminRepo: superAdmin,
	}
}

func (sa *SuperAdminHandler) CreateBranchHandler(e echo.Context) error {
	status, err := sa.SuperAdminRepo.CreateBranch(e)
	if err != nil {
		return e.JSON(status, echo.Map{
			"error": err.Error(),
		})
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (sa *SuperAdminHandler) DeleteBranchHandler(e echo.Context) error {
	status, err := sa.SuperAdminRepo.DeleteBranch(e)
	if err != nil {
		return e.JSON(status, echo.Map{
			"error": err.Error(),
		})
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}
