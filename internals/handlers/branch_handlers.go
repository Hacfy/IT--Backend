package handlers

import (
	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type BranchHandler struct {
	BranchRepo models.BranchInterface
}

func NewBranchHandler(branch models.BranchInterface) *BranchHandler {
	return &BranchHandler{
		BranchRepo: branch,
	}
}

func (bh *BranchHandler) CreateDepartmentHandler(e echo.Context) error {
	status, err := bh.BranchRepo.CreateDepartment(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (bh *BranchHandler) UpdateDepartmentHeadHandler(e echo.Context) error {
	status, err := bh.BranchRepo.UpdateDepartmentHead(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (bh *BranchHandler) CreateWarehouseHandler(e echo.Context) error {
	status, err := bh.BranchRepo.CreateWarehouse(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (bh *BranchHandler) UpdateWarehouseHeadHandler(e echo.Context) error {
	status, err := bh.BranchRepo.UpdateWarehouseHead(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}
