package handlers

import (
	"math"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type WarehouseHandler struct {
	WarehouseRepo models.WarehouseInterface
}

func NewWarehouse_Handler(warehouse models.WarehouseInterface) *WarehouseHandler {
	return &WarehouseHandler{
		WarehouseRepo: warehouse,
	}
}

func (wh *WarehouseHandler) CreateComponentHandler(e echo.Context) error {
	status, token, err := wh.WarehouseRepo.CreateComponent(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
		"token":   token,
	})
}

func (wh *WarehouseHandler) DeleteComponentHandler(e echo.Context) error {
	status, err := wh.WarehouseRepo.DeleteComponent(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (wh *WarehouseHandler) AddComponentUnitsHandler(e echo.Context) error {
	status, err := wh.WarehouseRepo.AddComponentUnits(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (wh *WarehouseHandler) AssignUnitsHandler(e echo.Context) error {
	status, err := wh.WarehouseRepo.AssignUnits(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (wh *WarehouseHandler) GetAllIssuesHandler(e echo.Context) error {
	status, issues, total, page, limit, err := wh.WarehouseRepo.GetAllWarehouseIssues(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"issues": issues,
		"meta": echo.Map{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func (wh *WarehouseHandler) GetAllWarehouseComponentsHandler(e echo.Context) error {
	status, components, err := wh.WarehouseRepo.GetAllWarehouseComponents(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}

	return e.JSON(status, echo.Map{
		"components": components,
		"total":      len(components),
	})
}
