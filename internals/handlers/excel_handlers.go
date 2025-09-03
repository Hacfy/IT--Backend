package handlers

import (
	"net/http"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/labstack/echo/v4"
)

type ExcelHandler struct {
	ExcelRepo models.ExcelInterface
}

func NewExcelHandler(excelRepo models.ExcelInterface) *ExcelHandler {
	return &ExcelHandler{
		ExcelRepo: excelRepo,
	}
}

func (eh *ExcelHandler) DownloadComponentMaintainanceReportHandler(e echo.Context) error {
	Status, File, err := eh.ExcelRepo.DownloadComponentMaintainanceReport(e)
	if err != nil {
		return echo.NewHTTPError(Status, err)
	}

	e.Response().Header().Set(echo.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	e.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=ComponentMaintenanceReport.xlsx")
	e.Response().WriteHeader(Status)

	if err := File.Write(e.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}
