package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) CreateSuperAdmin(superAdmin models.SuperAdminModel) (int, error) {
	var sa_id int

	var err error

	query1 := "INSERT INTO users(user_email, user_level) VALUES($1, $2)"
	query2 := "INSERT INTO super_admins(org_id, name, email, password) VALUES($1, $2, $3, $4) RETURNING id"

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

	if _, err = tx.Exec(query1, superAdmin.SuperAdminEmail, "super_admins"); err != nil {
		return -1, err
	}

	if err = tx.QueryRow(query2, superAdmin.Org_ID, superAdmin.SuperAdminName, superAdmin.SuperAdminEmail, superAdmin.SuperAdminPassword).Scan(&sa_id); err != nil {
		return -1, err
	}

	return sa_id, nil

}

func (q *Query) DeleteSuperAdmin(superAdminEmail string) (int, error) {
	var superAdminID, supersuperAdminOrgID int
	query0 := "SELECT EXISTS(SELECT 1 FROM branches WHERE super_admin_id = $1)"
	query1 := "DELETE FROM super_admins WHERE email = $1 RETURNING org_id, id"
	query3 := "DELETE FROM users WHERE user_email = $1"
	query2 := "INSERT INTO deleted_super_admins(super_admin_id, org_id, email) VALUES($1, $2, $3)"

	var err error

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

	var exists bool

	if err = tx.QueryRow(query0, superAdminID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if exists {
		return http.StatusConflict, fmt.Errorf("super_admin has branches associated with it")
	}

	if err = tx.QueryRow(query1, superAdminEmail).Scan(&supersuperAdminOrgID, &superAdminID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if _, err = tx.Exec(query2, superAdminID, supersuperAdminOrgID, superAdminEmail); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if _, err = tx.Exec(query3, superAdminEmail); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusNoContent, nil
}

func (q *Query) GetAllSuperAdmins(organisation_id int) ([]models.AllSuperAdminsDetailsModel, error) {
	query1 := "SELECT id, name, email FROM super_admins WHERE org_id = $1"
	query2 := "SELECT COUNT(*) FROM branches WHERE super_admin_id = $1"

	var err error

	var superAdmins []models.AllSuperAdminsDetailsModel
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

	rows, err = tx.Query(query1, organisation_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no superAdmins found for organisation %v", organisation_id)
			return nil, fmt.Errorf("no superAdmins found for organisation %v", organisation_id)
		}
		log.Printf("error while querying data: %v", err)
		return nil, fmt.Errorf("error while querying data")
	}
	defer rows.Close()

	for rows.Next() {
		var superAdmin models.AllSuperAdminsDetailsModel
		if err = rows.Scan(&superAdmin.SuperAdminID, &superAdmin.SuperAdminName, &superAdmin.SuperAdminEmail); err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error while scanning data")
		}

		var superAdminBranches int
		err = tx.QueryRow(query2, superAdmin.SuperAdminID).Scan(&superAdminBranches)
		if err != nil {
			log.Printf("error while scanning data: %v", err)
			return nil, fmt.Errorf("error while scanning data")
		}

		superAdmin.NoOfBranches = superAdminBranches
		superAdmins = append(superAdmins, superAdmin)
	}

	return superAdmins, nil
}

func (q *Query) ReassignSuperAdmin(superAdmin models.ReassignSuperAdminModel, org_id int) (int, error) {
	query := "UPDATE branches SET super_admin_id = $1 WHERE org_id = $2 AND super_admin_id = $3"
	if _, err := q.db.Exec(query, superAdmin.NewSuperAdminID, org_id, superAdmin.OldSuperAdminID); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

//query to get all the branches of a perticular organisation
// SELECT
//     d.department_id,
//     d.department_name,
//     COALESCE(dh.name, '') AS department_head_name,
//     COUNT(DISTINCT w.id) AS no_of_workspaces,
//     COUNT(DISTINCT i.id) AS issues
// FROM
//     departments d
// JOIN branches b ON d.branch_id = b.branch_id
// JOIN organisations o ON b.org_id = o.id
// LEFT JOIN department_heads dh ON d.department_id = dh.department_id
// LEFT JOIN workspaces w ON d.department_id = w.department_id
// LEFT JOIN issues i ON d.department_id = i.department_id
// WHERE
//     o.id = $1
// GROUP BY
//     d.department_id, d.department_name, dh.name
// ORDER BY
//     d.department_id;
