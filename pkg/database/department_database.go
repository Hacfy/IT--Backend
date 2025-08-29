package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) CreateWorkspace(workspace models.CreateWorkspaceModel, departmentHeadID int) (int, int, error) {
	query1 := "SELECT department_id FROM department_heads WHERE id = $1"
	query2 := "INSERT INTO workspaces(department_id, workspace_name) VALUES($1, $2) RETURNING id"

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	var department_id, workspace_id int

	if err := tx.QueryRow(query1, departmentHeadID).Scan(&department_id); err != nil {
		return http.StatusInternalServerError, -1, err
	}

	if err := tx.QueryRow(query2, department_id, workspace.WorkspaceName).Scan(&workspace_id); err != nil {
		return http.StatusInternalServerError, -1, err
	}

	return http.StatusCreated, workspace_id, nil
}

func (q *Query) DeleteWorkspace(workspace models.DeleteWorkspaceModel, departmentHeadID int) (int, error) {
	query1 := "DELETE FROM workspaces WHERE workspace_name = $1 AND workspace_id = $2"
	query2 := "INSERT INTO deleted_workspaces(workspace_id, department_id, deleted_by) VALUES($1, $2, $3)"

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

	var workspace_id int

	if _, err := tx.Exec(query1, workspace.WorkspaceName, workspace.WorkspaceID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query2, workspace_id, departmentHeadID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

func (q *Query) RaiseIssue(issue models.IssueModel) (int, int, error) {
	query1 := fmt.Sprintf("SELECT workspace_id FROM %s_units_assigned WHERE unit_id = $1", issue.UnitPrefix)
	query2 := `INSERT INTO issues (department_id, warehouse_id, workspace_id, unit_id, issue) VALUES($1, $2, $3, $4, $5) RETURNING id`
	query3 := fmt.Sprintf("UPDATE %s_us_units SET status = repair WHERE id = $1", issue.UnitPrefix)

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	var issue_id int
	var workspace_id int

	if err := tx.QueryRow(query1, issue.UnitID).Scan(&workspace_id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching unit found :%v", err)
			return http.StatusNotFound, -1, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, -1, err
	}

	if err := tx.QueryRow(query2, issue.DepartmentID, issue.WarehouseID, issue.WorkspaceID, issue.UnitID, issue.Issue).Scan(&issue_id); err != nil {
		return http.StatusInternalServerError, -1, err
	}

	if _, err := tx.Exec(query3, issue.UnitID); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching unit found :%v", err)
			return http.StatusNotFound, -1, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, -1, err
	}

	return http.StatusCreated, issue_id, nil
}

func (q *Query) RequestNewUnits(department_id int, workspace_id int, warehouse_id int, component_id int, number_of_units int, prefix string, user_id int) (int, int, error) {
	query1 := fmt.Sprintf("SELECT id FROM %s_units WHERE component_id = $1 AND warehouse_id = $2 AND id NOT IN (SELECT unit_id FROM %s_units_assigned )", prefix, prefix)
	query2 := "INSERT INTO requests(department_id, workspace_id, warehouse_id, component_id, number_of_units, prefix, created_by) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	var Request_id int
	var num int

	if err := tx.QueryRow(query1, component_id, warehouse_id).Scan(&num); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching unit found :%v", err)
			return http.StatusNotFound, -1, fmt.Errorf("no matching data found")
		}
		log.Printf("error while getting number of units: %v", err)
		return http.StatusInternalServerError, -1, fmt.Errorf("database error")
	}

	if num < number_of_units {
		log.Printf("not enough units available")
		return http.StatusBadRequest, -1, fmt.Errorf("not enough units available")
	}

	if err := tx.QueryRow(query2, department_id, workspace_id, warehouse_id, component_id, number_of_units, prefix, user_id).Scan(&Request_id); err != nil {
		log.Printf("error while requesting new units: %v", err)
		return http.StatusInternalServerError, -1, fmt.Errorf("database error")
	}

	return http.StatusCreated, Request_id, nil
}

func (q *Query) GetAllRequests(department_id int) ([]models.AllRequestsModel, error) {
	query := "SELECT id, workspace_id, warehouse_id, component_id, number_of_units, prefix, created_at, status FROM requests WHERE department_id = $1"

	var requests []models.AllRequestsModel

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

	rows, err := tx.Query(query, department_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching data found : %v", err)
			return nil, fmt.Errorf("no matching data found")
		}
		log.Printf("error while querying data: %v", err)
		return nil, fmt.Errorf("error occured while retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var request models.AllRequestsModel
		if err := rows.Scan(&request.RequestID, &request.WorkspaceID, &request.WarehouseID, &request.ComponentID, &request.NumberOfUnits, &request.Prefix, &request.CreatedAt, &request.Status); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error occured while retrieving data")
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
		return nil, fmt.Errorf("internal server error, please try again later")
	}

	return requests, nil
}

func (q *Query) GetRequestDetails(getRequestDetails models.GetRequestDetailsModel) (models.RequestDetailsModel, error) {
	query := "SELECT id, workspace_id, warehouse_id, component_id, number_of_units, prefix, created_by, created_at, status FROM requests WHERE id = $1"

	var Request models.RequestDetailsModel

	err := q.db.QueryRow(query, getRequestDetails.RequestID).Scan(&Request.RequestID, &Request.WorkspaceID, &Request.WarehouseID, &Request.ComponentID, &Request.NumberOfUnits, &Request.Prefix, &Request.CreatedBy, &Request.CreatedAt, &Request.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching data found : %v", err)
			return models.RequestDetailsModel{}, fmt.Errorf("no matching data found")
		}
		log.Printf("error while querying data: %v", err)
		return models.RequestDetailsModel{}, fmt.Errorf("error occured while retrieving data")
	}

	return Request, nil
}
