package handlers

import (
	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type DepartmentHandler struct {
	DepartmentRepo models.DepartmentInterface
}

func NewDepartmentHandler(department models.DepartmentInterface) *DepartmentHandler {
	return &DepartmentHandler{
		DepartmentRepo: department,
	}
}

func (dh *DepartmentHandler) CreateWorkspaceHandler(e echo.Context) error {
	status, workspace_id, err := dh.DepartmentRepo.CreateWorkspace(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}
	return e.JSON(status, echo.Map{
		"message":      "successfull",
		"workspace_id": workspace_id,
	})
}

func (dh *DepartmentHandler) DeleteWorkspaceHandler(e echo.Context) error {
	status, err := dh.DepartmentRepo.DeleteWorkspace(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
	})
}

func (dh *DepartmentHandler) RaiseIssueHandler(e echo.Context) error {
	status, issueID, err := dh.DepartmentRepo.RaiseIssue(e)
	if err != nil {
		return echo.NewHTTPError(status, err.Error())
	}
	return e.JSON(status, echo.Map{
		"message": "successfull",
		"issueID": issueID,
	})
}
