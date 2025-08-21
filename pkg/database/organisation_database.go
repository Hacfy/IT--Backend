package database

import (
	"log"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) CreateSuperAdmin(superAdmin models.SuperAdminModel) (int, error) {
	var sa_id int

	query1 := "INSERT INTO users(user_email, user_level) VALUES($1, $2)"
	query2 := "INSERT INTO super_admins(org_id, name, email, password) VALUES($1, $2, $3, $4) RETRUNING id"

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

	if _, err := tx.Exec(query1, superAdmin.SuperAdminEmail, "super_admins"); err != nil {
		return -1, err
	}

	if err := tx.QueryRow(query2, superAdmin.Org_ID, superAdmin.SuperAdminName, superAdmin.SuperAdminEmail, superAdmin.SuperAdminPassword).Scan(&sa_id); err != nil {
		return -1, err
	}

	return sa_id, nil

}

func (q *Query) DeleteSuperAdmin(superAdminEmail string) error {
	var superAdminID, supersuperAdminOrgID int
	query1 := "DELETE FROM super_admins WHERE email = $1 RETURNING org_id, id"
	query3 := "DELETE FROM users WHERE user_email = $1"
	query2 := "INSERT INTO deleted_super_admins(super_admin_id, org_id, email) VALUES($1, $2, $3)"

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	if err := tx.QueryRow(query1, superAdminEmail).Scan(&supersuperAdminOrgID, &superAdminID); err != nil {
		return err
	}

	if _, err := tx.Exec(query2, superAdminID, supersuperAdminOrgID, superAdminEmail); err != nil {
		return err
	}

	if _, err := tx.Exec(query3, superAdminEmail); err != nil {
		return err
	}

	return nil
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
