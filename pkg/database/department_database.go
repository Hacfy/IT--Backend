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
	query3 := fmt.Sprintf("UPDATE %s_us_units SET status = repair WHERE id = $1")

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
