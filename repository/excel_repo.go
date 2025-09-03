package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

type ExcelRepo struct {
	DB *sql.DB
}

func NewExcelRepo(db *sql.DB) *ExcelRepo {
	return &ExcelRepo{DB: db}
}
func (r *ExcelRepo) DownloadComponentMaintainanceReport(e echo.Context) (int, *excelize.File, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", r.DB)
	if err != nil {
		return status, nil, err
	}

	query := database.NewDBinstance(r.DB)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, nil, fmt.Errorf("invalid user details")
	}

	var request models.DownloadComponentMaintainanceReportRequest

	if err := e.Bind(&request); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, nil, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(request); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, nil, fmt.Errorf("failed to validate request")
	}

	ok, err = query.CheckIfComponentBelongsToWarehouse(request.ComponentID, claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, nil, fmt.Errorf("invalid user details")
	}

	componentName, componentPrefix, err := query.GetComponentNameAndPrefix(request.ComponentID)
	if err != nil {
		log.Printf("error while fetching component details: %v", err)
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	warehouseID, err := query.GetWarehouseIdOfComponent(request.ComponentID)
	if err != nil {
		log.Printf("error while fetching warehouse id: %v", err)
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	file := excelize.NewFile()
	sheet := "Component Maintenance Report"
	_, err = file.NewSheet(sheet)
	if err != nil {
		log.Printf("error while creating excel sheet: %v", err)
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	if err := file.MergeCell(sheet, "A1", "H1"); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}
	if err := file.SetCellValue(sheet, "A1", "Maintenance Report"); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	if err := file.MergeCell(sheet, "A2", "H2"); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}
	if err := file.SetCellValue(sheet, "A2", fmt.Sprintf("Warehouse ID: %d", warehouseID)); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	if err := file.MergeCell(sheet, "A3", "D3"); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}
	if err := file.SetCellValue(sheet, "A3", fmt.Sprintf("Component Name: %s", componentName)); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	if err := file.MergeCell(sheet, "E3", "H3"); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}
	if err := file.SetCellValue(sheet, "E3", fmt.Sprintf("Component Prefix: %s", componentPrefix)); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	if err := file.MergeCell(sheet, "A4", "D4"); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}
	if err := file.SetCellValue(sheet, "A4", fmt.Sprintf("Component ID: %d", request.ComponentID)); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	now := time.Now()
	today := now.Format("02-01-2006")

	if err := file.MergeCell(sheet, "E4", "H4"); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}
	if err := file.SetCellValue(sheet, "E4", today); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	headers := []string{"Unit ID",
		"Cost",
		"Maintenance Cost",
		"Last Maintenance Date",
		"Status",
		"Warranty",
		"Department ID",
		"Workspace ID",
	}

	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 5)
		if err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("database error")
		}
		if err := file.SetCellValue(sheet, cell, header); err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("database error")
		}
	}

	units, err := query.GetAllComponentUnits(componentPrefix)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("database error")
	}

	for i, unit := range units {
		row := i + 6
		_ = file.SetCellValue(sheet, fmt.Sprintf("A%d", row), unit.UnitID)
		_ = file.SetCellValue(sheet, fmt.Sprintf("B%d", row), unit.Cost)
		_ = file.SetCellValue(sheet, fmt.Sprintf("C%d", row), unit.MaintenanceCost)
		_ = file.SetCellValue(sheet, fmt.Sprintf("D%d", row), time.Unix(unit.LastMaintenanceDate, 0).Format("02-01-2006"))
		_ = file.SetCellValue(sheet, fmt.Sprintf("E%d", row), unit.Status)
		_ = file.SetCellValue(sheet, fmt.Sprintf("F%d", row), time.Unix(unit.WarrantyDate, 0).Format("02-01-2006"))
		_ = file.SetCellValue(sheet, fmt.Sprintf("G%d", row), unit.DepartmentID)
		_ = file.SetCellValue(sheet, fmt.Sprintf("H%d", row), unit.WorkspaceID)
	}

	titleStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 18, Color: "#000000"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	infoStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "#000000"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#D9D9D9"}},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	headerStyle, _ := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#000000"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	bodyStyle, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	_ = file.SetCellStyle(sheet, "A1", "H1", titleStyle)
	_ = file.SetCellStyle(sheet, "A2", "H2", infoStyle)
	_ = file.SetCellStyle(sheet, "A3", "D3", infoStyle)
	_ = file.SetCellStyle(sheet, "E3", "H3", infoStyle)
	_ = file.SetCellStyle(sheet, "A4", "D4", infoStyle)
	_ = file.SetCellStyle(sheet, "E4", "H4", infoStyle)
	_ = file.SetCellStyle(sheet, "A5", "H5", headerStyle)

	lastRow := len(units) + 5
	_ = file.SetCellStyle(sheet, "A6", fmt.Sprintf("H%d", lastRow), bodyStyle)

	for col := 1; col <= 8; col++ {
		maxLen := 0
		for row := 1; row <= lastRow; row++ {
			cell, _ := excelize.CoordinatesToCellName(col, row)
			val, _ := file.GetCellValue(sheet, cell)
			if len(val) > maxLen {
				maxLen = len(val)
			}
		}
		width := float64(maxLen) + 2
		colName, _ := excelize.ColumnNumberToName(col)
		_ = file.SetColWidth(sheet, colName, colName, width)
	}

	return http.StatusOK, file, nil
}
