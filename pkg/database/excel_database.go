package database

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) GetAllComponentUnits(prefix string) ([]models.ExcelMaintenanceReportModel, error) {
	query := "SELECT id, warranty_date, status, cost, maintainance_cost, last_maintenance_date FROM %s_units "

	var units []models.ExcelMaintenanceReportModel

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	var rows *sql.Rows

	rows, err = tx.Query(query, prefix)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no units found for component id %v", prefix)
			return nil, err
		}
		log.Printf("error while querying data: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var unit models.ExcelMaintenanceReportModel
		if err = rows.Scan(&unit.UnitID, &unit.WarrantyDate, &unit.Status, &unit.Cost, &unit.MaintenanceCost, &unit.LastMaintenanceDate); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, err
		}
		query1 := "SELECT department_id, workspace_id FROM %s_units_assigned WHERE unit_id = $1"
		var department_id, workspace_id int
		if err = tx.QueryRow(query1, prefix).Scan(&department_id, &workspace_id); err != nil {
			log.Printf("unit not assigned to any department or workspace: %v", unit.UnitID)
			unit.DepartmentID = "N/A"
			unit.WorkspaceID = "N/A"
		} else {
			department_id_str := strconv.Itoa(department_id)
			workspace_id_str := strconv.Itoa(workspace_id)
			unit.DepartmentID = department_id_str
			unit.WorkspaceID = workspace_id_str
		}
		units = append(units, unit)
	}

	if err = rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
		return nil, err
	}

	return units, nil
}

func (q *Query) GetAllComponentsPrefix(warehouse_id int) ([]models.ExcelPrefixReportModel, error) {
	query := "SELECT name, id, prefix FROM components WHERE warehouse_id = $1"

	var components []models.ExcelPrefixReportModel

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	var rows *sql.Rows

	rows, err = tx.Query(query, warehouse_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no components found for warehouse %v", warehouse_id)
			return nil, err
		}
		log.Printf("error while querying data: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var component models.ExcelPrefixReportModel
		if err = rows.Scan(&component.ComponentName, &component.ComponentID, &component.Prefix); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, err
		}
		components = append(components, component)
	}

	if err = rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
		return nil, err
	}

	return components, nil
}
