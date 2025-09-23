package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) CreateMainAdmin(main_admin models.MainAdminModel) (int, error) {
	var mainAdminID int
	query := "INSERT INTO main_admin(main_admin_email, main_admin_password) VALUES($1, $2) RETURNING main_admin_id"
	if err := q.db.QueryRow(query, main_admin.MainAdminEmail, main_admin.MainAdminPassword).Scan(&mainAdminID); err != nil {
		return -1, err
	}
	return mainAdminID, nil
}

func (q *Query) VerifyMainAdmin(main_admin_email string, main_admin_id int) (bool, error) {
	var exists int
	query := "SELECT 1 FROM main_admin WHERE main_admin_email = $1 AND main_admin_id = $2"
	if err := q.db.QueryRow(query, main_admin_email, main_admin_id).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (q *Query) Createorganization(organization models.OrganizationModel) (int, error) {
	var orgID int

	query1 := "INSERT INTO users(user_email, user_level) VALUES($1, $2)"
	query2 := "INSERT INTO organization(main_admin_id, name, email, phone_number, password) VALUES($1, $2, $3, $4, $5) RETURNING id"

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

	if _, err = tx.Exec(query1, organization.OrganizationEmail, "organization"); err != nil {
		return -1, err
	}

	if err = tx.QueryRow(query2, organization.OrganizationMainAdminID, organization.OrganizationName, organization.OrganizationEmail, organization.OrganizationPhoneNumber, organization.OrganizationPassword).Scan(&orgID); err != nil {
		return -1, err
	}

	return orgID, nil
}

func (q *Query) GetMainAdminCredentials(main_admin_email string) (models.MainAdminModel, bool, error) {
	var main_admin_id int
	var main_admin_password string
	query := "SELECT main_admin_id, main_admin_password FROM main_admin WHERE main_admin_email = $1"
	if err := q.db.QueryRow(query, main_admin_email).Scan(&main_admin_id, &main_admin_password); err != nil {
		if err == sql.ErrNoRows {
			return models.MainAdminModel{}, false, nil
		}
		return models.MainAdminModel{}, false, err
	}
	return models.MainAdminModel{
		MainAdminID:       main_admin_id,
		MainAdminEmail:    main_admin_email,
		MainAdminPassword: main_admin_password,
	}, true, nil
}

func (q *Query) DeleteMainAdmin(mainAdminEmail string, main_admin_id, deleted_by int) (int, error) {
	query0 := "SELECT EXISTS(SELECT 1 FROM organization "
	query1 := "DELETE FROM main_admin WHERE main_admin_email = $1"
	query2 := "INSERT INTO deleted_main_admins(main_admin_id, email, deleted_by) VALUES($1, $2, $3)"

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

	if err = tx.QueryRow(query0, main_admin_id).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if exists {
		return http.StatusConflict, fmt.Errorf("main_admin has organization associated with them")
	}

	if _, err = tx.Exec(query1, mainAdminEmail); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if _, err = tx.Exec(query2, main_admin_id, mainAdminEmail, deleted_by); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusNoContent, nil

}

func (q *Query) Deleteorganization(organizationEmail string, organization_id, deleted_by int) (int, error) {
	query1 := "SELECT EXISTS(SELECT 1 FROM super_admins WHERE org_id = $1)"
	query2 := "INSERT INTO deleted_organization(org_id, email, main_admin_id) VALUES($1, $2, $3)"
	query3 := "DELETE FROM organization WHERE email = $1"
	query4 := "DELETE FROM users WHERE user_email = $1 RETRUNING user_email, user_level, ever_logged_in, latest_token, created_at"
	query5 := "INSERT INTO deleted_users(user_email, user_level, ever_logged_in, latest_token, created_at, deleted_by)"

	var org_email string
	var org_level string
	var super_admin_exists bool
	var ever_logged_in bool
	var latest_token string
	var created_at int64

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

	if err = tx.QueryRow(query1, organizationEmail).Scan(&super_admin_exists); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if super_admin_exists {
		return http.StatusConflict, fmt.Errorf("super_admin has branches associated with it")
	}
	if err = tx.QueryRow(query2, organization_id, organizationEmail, deleted_by).Scan(&org_email, &org_level, &ever_logged_in, &latest_token, &created_at); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if _, err = tx.Exec(query3, organizationEmail); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if err = tx.QueryRow(query4, organization_id, organizationEmail, deleted_by).Scan(&org_email, &org_level, &ever_logged_in, &latest_token, &created_at); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if err = tx.QueryRow(query5, organization_id, organizationEmail, deleted_by).Scan(&org_email, &org_level, &ever_logged_in, &latest_token, &created_at); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusNoContent, nil
}

func (q *Query) GetAllorganization(mainAdminID int) ([]models.GetAllOrganizationModel, error) {
	var orgs []models.GetAllOrganizationModel
	query := "SELECT id, name, email, phone_number FROM organization WHERE main_admin_id = $1"

	rows, err := q.db.Query(query, mainAdminID)
	if err != nil {
		return []models.GetAllOrganizationModel{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var org models.GetAllOrganizationModel
		err = rows.Scan(&org.OrganizationID, &org.OrganizationName, &org.OrganizationEmail, &org.OrganizationPhoneNumber)
		if err != nil {
			return []models.GetAllOrganizationModel{}, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}

func (q *Query) GetAllMainAdmins() ([]models.AllMainAdminModel, error) {
	var main_admins []models.AllMainAdminModel

	query := "SELECT main_admin_id, main_admin_email, COUNT(organization_id) AS no_of_orgs FROM main_admin JOIN organization ON main_admin_id = main_admin_id"

	rows, err := q.db.Query(query)
	if err != nil {
		return []models.AllMainAdminModel{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var main_admin models.AllMainAdminModel
		err := rows.Scan(&main_admin.MainAdminID, &main_admin.MainAdminEmail, &main_admin.NoOfOrgs)
		if err != nil {
			return []models.AllMainAdminModel{}, err
		}
		main_admins = append(main_admins, main_admin)
	}

	return main_admins, nil
}
