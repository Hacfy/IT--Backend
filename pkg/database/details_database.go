package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) GetAllDepartments(branch_id int, sort models.SortModel) (int, []models.AllDepartmentsModel, int, error) {

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, []models.AllDepartmentsModel{}, -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	whereClause := fmt.Sprintf("WHERE d.branch_id = %d ", branch_id)

	if sort.Search != "" {
		id, err := strconv.Atoi(sort.Search)
		if err == nil {
			whereClause += fmt.Sprintf("AND (d.department_id = %d OR  dh.id = %d OR w.id = %d )", id, id, id)
		} else {
			whereClause += fmt.Sprintf("AND (d.department_name LIKE '%%%s%%' OR dh.name LIKE '%%%s%%' OR w.workspace_name LIKE '%%%s%%')", sort.Search, sort.Search, sort.Search)
		}
	}

	query1 := fmt.Sprintf(`SELECT 
		d.department_id,
		d.department_name,
		COALESCE(dh.name, '') AS department_head_name,
		COUNT(DISTINCT w.id) AS no_of_workspaces,
		COUNT(DISTINCT i.id) AS issues
	FROM 
		departments d
	LEFT JOIN department_heads dh ON d.department_id = dh.department_id
	LEFT JOIN workspaces w ON d.department_id = w.department_id
	LEFT JOIN issues i ON d.department_id = i.department_id
	%s
	GROUP BY 
		d.department_id, d.department_name, dh.name
	ORDER BY 
		%s %s
	LIMIT $1 OFFSET $2;
		`, whereClause, sort.SortBy, sort.Order)
	query2 := `SELECT COUNT(*) AS total_departments 
		FROM departments d
		WHERE branch_id = $1`

	rows, err := tx.Query(query1, sort.Limit, sort.Offset)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no departments with branch_id %v found: %v", branch_id, err)
			return http.StatusNotFound, []models.AllDepartmentsModel{}, 0, fmt.Errorf("no departments found")
		}
		log.Printf("error while getting department data: %v", err)
		return http.StatusInternalServerError, []models.AllDepartmentsModel{}, -1, fmt.Errorf("internal server error, please try again later")
	}
	defer rows.Close()

	var Departments []models.AllDepartmentsModel

	for rows.Next() {
		var dept models.AllDepartmentsModel
		if err := rows.Scan(
			&dept.DepartmentID,
			&dept.DepartmentName,
			&dept.DepartmentHeadName,
			&dept.NoOfWorkspaces,
			&dept.Issues,
		); err != nil {
			log.Printf("error while row iteration: %v", err)
			return http.StatusInternalServerError, []models.AllDepartmentsModel{}, -1, fmt.Errorf("internal server error, please try again later")
		}
		Departments = append(Departments, dept)
	}

	var Total_departments int

	if err = tx.QueryRow(query2, branch_id).Scan(&Total_departments); err != nil {
		log.Printf("errro while getting total no. of departments in branch %v: %v", branch_id, err)
		return http.StatusInternalServerError, []models.AllDepartmentsModel{}, -1, fmt.Errorf("internal server error, please try again later")
	}

	return http.StatusOK, Departments, Total_departments, nil

}

func (q *Query) GetDepartmentIssues(department_id int, sort models.SortModel) (int, []models.DepartmentIssuesModel, int, error) {

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, nil, -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	whereClause := fmt.Sprintf("WHERE d.department_id = %d ", department_id)

	if sort.Search != "" {
		id, err := strconv.Atoi(sort.Search)
		if err == nil {
			whereClause += fmt.Sprintf("AND (d.department_id = %d OR  i.id = %d OR w.id = %d )", id, id, id)
		} else {
			whereClause += fmt.Sprintf("AND (i.issue LIKE '%%%s%%' OR w.workspace_name LIKE '%%%s%%')", sort.Search, sort.Search)
		}
	}

	query1 := fmt.Sprintf(`SELECT
	i.id, 
	i.issue, 
	i.created_at, 
	i.status, 
	i.unit_id,
	i.unit_prefix,
	w.workspace_id
	FROM 
		issues i
	LEFT JOIN workspaces w ON i.workspace_id = w.id
	LEFT JOIN departments d ON i.department_id = d.department_id
	LEFT JOIN department_heads dh ON d.department_id = dh.department_id
	%s
	ORDER BY i.%s %s
	LIMIT $1 OFFSET $2;`,
		whereClause, sort.SortBy, sort.Order)

	query2 := `SELECT COUNT(*) AS total_issues 
		FROM issues i
		WHERE i.department_id = $1`

	rows, err := tx.Query(query1, sort.Limit, sort.Offset)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no issues with department_id %v found: %v", department_id, err)
			return http.StatusNotFound, nil, 0, fmt.Errorf("no issues found")
		}
		log.Printf("error while getting issue data: %v", err)
		return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
	}
	defer rows.Close()

	var Issues []models.DepartmentIssuesModel

	for rows.Next() {
		var issue models.DepartmentIssuesModel
		if err := rows.Scan(
			&issue.IssueID,
			&issue.Issue,
			&issue.CreatedAt,
			&issue.Status,
			&issue.UnitID,
			&issue.UnitPrefix,
			&issue.WorkspaceID,
		); err != nil {
			log.Printf("error while row iteration: %v", err)
			return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
		}
		Issues = append(Issues, issue)
	}

	var Total_issues int

	if err = tx.QueryRow(query2, department_id).Scan(&Total_issues); err != nil {
		log.Printf("errro while getting total no. of issues in department %v: %v", department_id, err)
		return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
	}

	return http.StatusOK, Issues, Total_issues, nil
}

func (q *Query) GetAllWorkspaces(department_id int, sort models.SortModel) (int, []models.DepartmentWorkspaceModel, int, error) {

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, nil, -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	whereClause := fmt.Sprintf("WHERE w.department_id = %d ", department_id)

	if sort.Search != "" {
		id, err := strconv.Atoi(sort.Search)
		if err == nil {
			whereClause += fmt.Sprintf("AND (w.id = %d OR  w.department_id = %d )", id, id)
		} else {
			whereClause += fmt.Sprintf("AND (w.workspace_name LIKE '%%%s%%' )", sort.Search)
		}
	}

	query1 := fmt.Sprintf(`
			SELECT
	        w.id,
	        w.workspace_name
			FROM
				workspaces w
			%s
			ORDER BY w.%s %s
			LIMIT $1 OFFSET $2;
		`,
		whereClause, sort.SortBy, sort.Order)

	query2 := `SELECT COUNT(*) AS total_workspaces
		FROM workspaces w
		WHERE w.department_id = $1`

	rows, err := tx.Query(query1, sort.Limit, sort.Offset)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no workspaces with department_id %v found: %v", department_id, err)
			return http.StatusNotFound, nil, 0, fmt.Errorf("no workspaces found")
		}
		log.Printf("error while getting workspace data: %v", err)
		return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
	}
	defer rows.Close()

	var Workspaces []models.DepartmentWorkspaceModel
	for rows.Next() {
		var workspace models.DepartmentWorkspaceModel
		if err := rows.Scan(
			&workspace.WorkspaceID,
			&workspace.WorkspaceName,
		); err != nil {
			log.Printf("error while row iteration: %v", err)
			return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
		}
		Workspaces = append(Workspaces, workspace)
	}

	var Total_workspaces int

	if err = tx.QueryRow(query2, department_id).Scan(&Total_workspaces); err != nil {
		log.Printf("errro while getting total no. of workspaces in department %v: %v", department_id, err)
		return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
	}

	return http.StatusOK, Workspaces, Total_workspaces, nil
}

func (q *Query) GetAllBranches(super_admin_id int, sort models.SortModel) (int, []models.AllBranchesModel, int, error) {

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return http.StatusInternalServerError, nil, -1, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	whereClause := fmt.Sprintf("WHERE b.super_admin_id = %d ", super_admin_id)

	if sort.Search != "" {
		id, err := strconv.Atoi(sort.Search)
		if err == nil {
			whereClause += fmt.Sprintf("AND (b.branch_id = %d )", id)
		} else {
			whereClause += fmt.Sprintf("AND (b.branch_location LIKE '%%%s%%' OR b.branch_name LIKE '%%%s%%')", sort.Search, sort.Search)
		}
	}

	query1 := fmt.Sprintf(`
			SELECT
				b.branch_id,
				b.branch_name,
				b.branch_location,
				bh.name AS branch_head_name
			FROM
				branches b
			LEFT JOIN branch_heads bh ON b.branch_id = bh.branch_id
			%s
			ORDER BY b.%s %s
			LIMIT $1 OFFSET $2;
		`,
		whereClause, sort.SortBy, sort.Order)

	query2 := `SELECT COUNT(*) AS total_branches
			FROM branches b
			WHERE b.super_admin_id = $1`

	rows, err := tx.Query(query1, sort.Limit, sort.Offset)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no branches with super_admin_id %v found: %v", super_admin_id, err)
			return http.StatusNotFound, nil, 0, fmt.Errorf("no branches found")
		}
		log.Printf("error while getting branch data: %v", err)
		return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
	}
	defer rows.Close()

	var Branches []models.AllBranchesModel

	for rows.Next() {
		var branch models.AllBranchesModel
		if err := rows.Scan(
			&branch.BranchID,
			&branch.BranchName,
			&branch.BranchLocation,
			&branch.BranchHeadName,
		); err != nil {
			log.Printf("error while row iteration: %v", err)
			return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
		}
		Branches = append(Branches, branch)
	}

	var Total_branches int

	if err = tx.QueryRow(query2, super_admin_id).Scan(&Total_branches); err != nil {
		log.Printf("errro while getting total no. of branches in super_admin %v: %v", super_admin_id, err)
		return http.StatusInternalServerError, nil, -1, fmt.Errorf("internal server error, please try again later")
	}
	return http.StatusOK, Branches, Total_branches, nil
}

func (q *Query) CheckIfDepartmentUnderBranchHead(department_id, user_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM departments WHERE department_id = $1 AND branch_id = (SELECT branch_id FROM branch_heads WHERE id = $2)"
	var exists bool
	err := q.db.QueryRow(query, department_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (q *Query) CheckIfWarehouseUnderBranchHead(warehouse_id, user_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM warehouses WHERE warehouse_id = $1 AND branch_id = (SELECT branch_id FROM branch_heads WHERE id = $2)"
	var exists bool
	err := q.db.QueryRow(query, warehouse_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (q *Query) CheckIfWarehouseIDExistsInTheDepartmentsBranch(warehouse_id, user_id int) (bool, error) {
	query := `SELECT EXISTS (
				SELECT 1
				FROM warehouses w
				JOIN departments d ON d.department_id = dh.department_id
				JOIN department_heads dh ON dh.department_id = d.department_id
				WHERE w.id = $1 AND dh.id = $2 AND w.branch_id = d.branch_id
			) AS same_branch;`

	var same_branch bool

	err := q.db.QueryRow(query, warehouse_id, user_id).Scan(&same_branch)

	return same_branch, err
}

func (q *Query) CheckBranchHead(user_id, branch_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM branch_heads WHERE id = $1 AND branch_id = $2)"
	var exists bool
	err := q.db.QueryRow(query, user_id, branch_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (q *Query) CheckIfBranchUnderSuperAdmin(branch_id, user_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM branches WHERE branch_id = $1 AND super_admin_id = (SELECT id FROM super_admins WHERE id = $2))"
	var exists bool
	err := q.db.QueryRow(query, branch_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (q *Query) CheckIfBranchUnderOrganisationAdmin(branch_id, user_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM branches WHERE branch_id = $1 AND org_id = (SELECT id FROM organisations WHERE id = $2))"
	var exists bool
	err := q.db.QueryRow(query, branch_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (q *Query) CheckDepartmentHead(user_id, department_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM department_heads WHERE id = $1 AND department_id = $2)"
	var exists bool
	err := q.db.QueryRow(query, user_id, department_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (q *Query) CheckIfDepartmentUnderSuperAdmin(department_id, user_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM departments WHERE department_id = $1 AND branch_id IN (SELECT branch_id FROM branches WHERE super_admin_id = (SELECT id FROM super_admins WHERE id = $2)))"
	var exists bool
	err := q.db.QueryRow(query, department_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (q *Query) CheckIfDepartmentUnderOrganisationAdmin(department_id, user_id int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM departments WHERE department_id = $1 AND branch_id IN (SELECT branch_id FROM branches WHERE org_id = (SELECT id FROM organisations WHERE id = $2)))"
	var exists bool
	err := q.db.QueryRow(query, department_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
