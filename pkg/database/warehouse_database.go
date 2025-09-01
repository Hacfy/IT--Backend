package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) IfPrefixExists(prefix string) bool {
	var exists bool
	q.db.QueryRow("SELECT EXISTS(SELECT 1 FROM warehouses WHERE prefix = $1)", prefix).Scan(&exists)
	return exists
}

func (q *Query) IfComponentExists(name string, warehouse_id int) bool {
	var exists bool
	q.db.QueryRow("SELECT EXISTS(SELECT 1 FROM components WHERE name = $1 AND warehouse_id = $2)", name, warehouse_id).Scan(&exists)
	return exists
}

func (q *Query) CheckIfComponentIDExists(component_id, warehouse_id int) (string, bool, error) {
	query := "SELECT prefix FROM components WHERE id = $1 AND warehouse_id = $2"
	var prefix string
	err := q.db.QueryRow(query, component_id, warehouse_id).Scan(&prefix)
	if err != nil {
		return "", false, err
	}
	return prefix, true, nil
}

func (q *Query) CreateComponent(name, prefix string, warehouse_id int) (int, error) {
	var id int
	query1 := "INSERT INTO components(name, prefix, warehouse_id) VALUES($1, $2, $3) RETURNING id"
	query2 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s_units (
		id SERIAL PRIMARY KEY,
		component_id INTEGER NOT NULL,
		warehouse_id INTEGER NOT NULL,
		warranty_date TIMESTAMPTZ NOT NULL,
		status unit_status DEFAULT 'working',
		cost NUMERIC(10, 2) NOT NULL,
		maintainance_cost NUMERIC(10, 2) DEFAULT 0,
		CONSTRAINT fk_%s_units_component_id FOREIGN KEY (component_id) REFERENCES components(id) ON DELETE CASCADE,
		CONSTRAINT fk_%s_units_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE
	)`, prefix, prefix, prefix)
	query3 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s_units_assigned (
		id INTEGER PRIMARY KEY,
		department_id INTEGER NOT NULL,
		workspace_id INTEGER NOT NULL,
		assigned_at TIMESTAMPTZ DEFAULT NOW(),
		CONSTRAINT fk_%s_units_department_id FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
		CONSTRAINT fk_%s_units_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		CONSTRAINT fk_%s_units_id FOREIGN KEY (id) REFERENCES %s_units(id) ON DELETE CASCADE
	)`, prefix, prefix, prefix, prefix, prefix)

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	if err := tx.QueryRow(query1, name, prefix, warehouse_id).Scan(&id); err != nil {
		return -1, err
	}

	if _, err := tx.Exec(query2); err != nil {
		return -1, err
	}

	if _, err := tx.Exec(query3); err != nil {
		return -1, err
	}

	return id, nil
}

func (q *Query) DeleteComponent(del_component models.DeleteComponentModel, warehouse_id int) (int, error) {
	query1 := "CALL delete_component($1, $2)"

	if _, err := q.db.Exec(query1, del_component.ComponentID, warehouse_id); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}
	return http.StatusNoContent, nil
}

func (q *Query) CreateComponentUnit(warranty_date time.Time, cost float32, prifix string, warehouse_id, number, component_id int) (int, error) {
	query1 := fmt.Sprintf("INSERT INTO %s_units(component_id, warehouse_id, warranty_date, cost) VALUES($1, $2, $3, $4)", prifix)

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	for range number {
		_, err := tx.Exec(query1, component_id, warehouse_id, warranty_date, cost)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusCreated, nil

}

func (q *Query) AssignUnitWorkspace(workspace_id int, unit_id []int, prefix string) (int, error) {
	query1 := "SELECT department_id FROM workspaces WHERE id = $1"
	query2 := fmt.Sprintf("INSERT INTO %s_units_assigned(department_id, workspace_id, id) VALUES($1, $2, $3)", prefix)

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	var department_id int

	if err := tx.QueryRow(query1, workspace_id).Scan(&department_id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching department found")
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		log.Printf("error while getting department id: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	for _, unit := range unit_id {
		if _, err := tx.Exec(query2, department_id, workspace_id, unit); err != nil {
			log.Printf("error while assigning units: %v", err)
			return http.StatusInternalServerError, fmt.Errorf("database error")
		}
	}

	return http.StatusOK, nil

}

func (q *Query) GetAllIssues(id int, sort models.SortModel) (int, []models.IssueModel, int, error) {
	var issues []models.IssueModel

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, []models.IssueModel{}, 0, fmt.Errorf("database error")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	args := []interface{}{sort.Limit, sort.Offset}
	argIndex := 3
	whereClause := fmt.Sprintf("WHERE workspace_id = %d ", id)

	if sort.Search != "" {
		id, err := strconv.Atoi(sort.Search)
		if err == nil {
			whereClause += fmt.Sprintf("AND (id = $%d OR department_id = $%d OR workspace_id = $%d OR unit_id = $%d OR created_at = $%d) ", argIndex, argIndex, argIndex, argIndex, argIndex)
			args = append(args, id)
		} else {
			whereClause += fmt.Sprintf("AND (unit_prefix ILIKE $%d OR issue ILIKE $%d OR status ILIKE $%d) ", argIndex, argIndex, argIndex)
			args = append(args, "%"+sort.Search+"%")
		}
	}

	query := fmt.Sprintf(`SELECT id, department_id, workspace_id, unit_id, unit_prefix, issue, created_at, status FROM issues
							%s
							ORDER BY %s %s 
							LIMIT $1 OFFSET $2
							`, whereClause, sort.SortBy, sort.Order)

	rows, err := tx.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching data found : %v", err)
			return http.StatusNotFound, []models.IssueModel{}, 0, fmt.Errorf("no matching data found")
		}
		log.Printf("error while querying data: %v", err)
		return http.StatusInternalServerError, []models.IssueModel{}, 0, fmt.Errorf("error occured while retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var issue models.IssueModel
		if err := rows.Scan(&issue.IssueID, &issue.DepartmentID, &issue.WorkspaceID, &issue.UnitID, &issue.UnitPrefix, &issue.Issue, &issue.Created_at, &issue.Status); err != nil {
			log.Printf("error while scanning data: %v", err)
			return http.StatusInternalServerError, []models.IssueModel{}, 0, fmt.Errorf("error occured while retrieving data")
		}
		issues = append(issues, issue)
	}

	if err := rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
		return http.StatusInternalServerError, []models.IssueModel{}, 0, fmt.Errorf("internal server error, please try again later")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM issues %s`, whereClause)
	var total int
	if err := tx.QueryRow(countQuery).Scan(&total); err != nil {
		log.Printf("error while scanning data: %v", err)
		return http.StatusInternalServerError, []models.IssueModel{}, 0, fmt.Errorf("error occured while retrieving data")
	}

	return http.StatusOK, issues, total, nil

}

func (q *Query) GetAllWarehouseComponents(warehouse_id int) ([]models.AllWarehouseComponentsModel, error) {
	query := "SELECT id, name, prefix FROM components WHERE warehouse_id = $1"
	var components []models.AllWarehouseComponentsModel
	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return nil, fmt.Errorf("database error")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	rows, err := tx.Query(query, warehouse_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no components found for user %v", warehouse_id)
			return nil, fmt.Errorf("no components found")
		}
		log.Printf("error while querying data: %v", err)
		return nil, fmt.Errorf("error occured while retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var component models.AllWarehouseComponentsModel
		if err := rows.Scan(
			&component.ComponentID,
			&component.ComponentName,
			&component.Prefix,
		); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error occured while retrieving data")
		}

		Query := fmt.Sprintf(`SELECT COUNT(*) FROM %s_units WHERE component_id = $1`, component.Prefix)

		var units int

		if err := tx.QueryRow(Query, component.ComponentID).Scan(&units); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error occured while retrieving data")
		}

		component.Units = units
		components = append(components, component)
	}

	return components, nil
}

func (q *Query) GetAllWarehouseComponentUnits(component_id int) ([]models.AllComponentUnitsModel, error) {
	query1 := "SELECT prefix FROM components WHERE id = $1"

	var units []models.AllComponentUnitsModel

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return nil, fmt.Errorf("database error")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	var prefix string

	err = tx.QueryRow(query1, component_id).Scan(&prefix)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no components found for component id %v", component_id)
			return nil, fmt.Errorf("no components found")
		}
		log.Printf("error while querying data: %v", err)
		return nil, fmt.Errorf("error occured while retrieving data")
	}

	query2 := fmt.Sprintf(`SELECT id, warehouse_id, warranty_date, status, cost, maintainance_cost FROM %s_units WHERE component_id = $1`, prefix)
	rows, err := tx.Query(query2, component_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no units found for component id %v", component_id)
			return nil, fmt.Errorf("no units found")
		}
		log.Printf("error while querying data: %v", err)
		return nil, fmt.Errorf("error occured while retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var unit models.AllComponentUnitsModel
		if err := rows.Scan(&unit.UnitID, &unit.WarehouseID, &unit.WarrantyDate, &unit.Status, &unit.Cost, &unit.MaintenanceCost); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error occured while retrieving data")
		}
		query3 := fmt.Sprintf(`SELECT COUNT(*) FROM %s_units_assigned WHERE unit_id = $1`, prefix)

		var units_assigned int

		if err := tx.QueryRow(query3, unit.UnitID).Scan(&units_assigned); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error occured while retrieving data")
		}
		if units_assigned > 0 {
			unit.Assigned = true
		} else {
			unit.Assigned = false
		}
		units = append(units, unit)
	}

	return units, nil
}

func (q *Query) GetIssueDetails(issue_id int) (models.IssueDetailsModel, error) {
	query1 := "SELECT department_id, warehouse_id, workspace_id, unit_id, unit_prefix, issue, created_at, status FROM issues WHERE id = $1"

	var issue models.IssueDetailsModel

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return issue, fmt.Errorf("database error")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	err = tx.QueryRow(query1, issue_id).Scan(&issue.DepartmentID, &issue.WarehouseID, &issue.WorkspaceID, &issue.UnitID, &issue.UnitPrefix, &issue.Issue, &issue.Created_at, &issue.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching data found")
			return issue, fmt.Errorf("no matching data found")
		}
		log.Printf("error while querying data: %v", err)
		return issue, fmt.Errorf("error occured while retrieving data")
	}

	return issue, nil
}

func (q *Query) GetUnitAssignmentHistory(unit_id int) ([]models.HistoryModel, error) {
	query := "SELECT workspace_id, department_id, assigned_at, deleted_at, deleted_by FROM deleted_units_assigned WHERE unit_id = $1"

	var History []models.HistoryModel

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return nil, fmt.Errorf("database error")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	rows, err := tx.Query(query, unit_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching data found")
			return nil, fmt.Errorf("no matching data found")
		}
		log.Printf("error while querying data: %v", err)
		return nil, fmt.Errorf("error occured while retrieving data")
	}

	defer rows.Close()

	for rows.Next() {
		var history models.HistoryModel
		if err := rows.Scan(&history.WorkspaceID, &history.DepartmentID, &history.AssignedAt, &history.DeletedAt, &history.DeletedBy); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error occured while retrieving data")
		}
		History = append(History, history)
	}

	return History, nil
}

func (q *Query) UpdateIssueStatus(issue_id int, status string) (int, error) {
	query := "UPDATE issues SET status = $1 WHERE id = $2"

	if _, err := q.db.Exec(query, status, issue_id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching data found")
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		log.Printf("error while updating issue status: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusOK, nil
}

func (q *Query) UpdateComponentName(component_id int, component_name string) (int, error) {
	query := "UPDATE components SET name = $1 WHERE id = $2"

	if _, err := q.db.Exec(query, component_name, component_id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching data found")
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		log.Printf("error while updating component name: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusOK, nil
}
