package handlers

import (
	"math"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type DetailsHandler struct {
	DetailsRepo models.DetailsInterface
}

func NewDetailsHandler(detailsRepo models.DetailsInterface) *DetailsHandler {
	return &DetailsHandler{
		DetailsRepo: detailsRepo,
	}
}

func (dh *DetailsHandler) GetAllDepartmentsHandler(e echo.Context) error {
	Departments, Status, Total, Page, Limit, err := dh.DetailsRepo.GetAllDepartmentsRepo(e)
	if err != nil {
		return echo.NewHTTPError(Status, err)
	}

	return e.JSON(Status, echo.Map{
		"departments": Departments,
		"meta": echo.Map{
			"total": Total,
			"page":  Page,
			"limit": Limit,
			"pages": int(math.Ceil(float64(Total) / float64(Limit))),
		},
	})
}

func (dh *DetailsHandler) GetDepartmentIssuesHandler(e echo.Context) error {
	Status, Issues, Total, Page, Limit, err := dh.DetailsRepo.GetDepartmentIssues(e)
	if err != nil {
		return echo.NewHTTPError(Status, err)
	}

	return e.JSON(Status, echo.Map{
		"issues": Issues,
		"meta": echo.Map{
			"total": Total,
			"page":  Page,
			"limit": Limit,
			"pages": int(math.Ceil(float64(Total) / float64(Limit))),
		},
	})
}

func (dh *DetailsHandler) GetDepartmentWorkspacesHandler(e echo.Context) error {
	Workspaces, Status, Total, Page, Limit, err := dh.DetailsRepo.GetDepartmentWorkspaces(e)
	if err != nil {
		return echo.NewHTTPError(Status, err)
	}

	return e.JSON(Status, echo.Map{
		"workspaces": Workspaces,
		"meta": echo.Map{
			"total": Total,
			"page":  Page,
			"limit": Limit,
			"pages": int(math.Ceil(float64(Total) / float64(Limit))),
		},
	})
}

func (db *DetailsHandler) GetAllBranchesHandler(e echo.Context) error {
	Branches, Status, Total, Page, Limit, err := db.DetailsRepo.GetAllBranchesUnderSuperAdmin(e)
	if err != nil {
		return echo.NewHTTPError(Status, err)
	}

	return e.JSON(Status, echo.Map{
		"branches": Branches,
		"meta": echo.Map{
			"total": Total,
			"page":  Page,
			"limit": Limit,
			"pages": int(math.Ceil(float64(Total) / float64(Limit))),
		},
	})
}

func (db *DetailsHandler) GetAllWarehousesHandler(e echo.Context) error {
	Warehouses, Status, err := db.DetailsRepo.GetAllWarehouses(e)
	if err != nil {
		return echo.NewHTTPError(Status, err)
	}

	return e.JSON(Status, echo.Map{
		"warehouses": Warehouses,
		"total":      len(Warehouses),
	})
}

func (detailsHandler DetailsHandler) GetAllDepartmentOutOfWarentyUnitsHandler(e echo.Context) error {
	status, OutOfWarentyUnits, total, Limit, Page, err := detailsHandler.DetailsRepo.GetAllOutOfWarehouseUnitsInDepartment(e)
	if err != nil {
		return echo.NewHTTPError(status, err)
	}

	return e.JSON(status, echo.Map{
		"outOfWarentyUnits": OutOfWarentyUnits,
		"meta": echo.Map{
			"total": total,
			"page":  Page,
			"limit": Limit,
			"pages": int(math.Ceil(float64(total) / float64(Limit))),
		},
	})
}
