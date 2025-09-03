package models

import (
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

type DownloadComponentMaintainanceReportRequest struct {
	ComponentID int `json:"component_id" validate:"required"`
}

type ExcelMaintenanceReportModel struct {
	UnitID              int    `json:"unit_id"`
	WarrantyDate        int64  `json:"warranty_date"`
	Status              string `json:"status"`
	Cost                int    `json:"cost"`
	MaintenanceCost     int    `json:"maintainance_cost"`
	LastMaintenanceDate int64  `json:"last_maintenance_date"`
	DepartmentID        string `json:"department_id"`
	WorkspaceID         string `json:"workspace_id"`
}

type ExcelInterface interface {
	DownloadComponentMaintainanceReport(echo.Context) (int, *excelize.File, error)
}
