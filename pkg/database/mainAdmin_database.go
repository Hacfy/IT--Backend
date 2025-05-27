package database

import (
	"database/sql"
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

func (q *Query) CreateOrganisation(organisation models.OrganisationModel) (int, error) {
	var orgID int

	query1 := "INSERT INTO users(user_email, user_level) VALUES($1, $2)"
	query2 := "INSERT INTO organisations(main_admin_id, name, email, phone_number, password) VALUES($1, $2, $3, $4, $5) RETURNING id"

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

	if _, err := tx.Exec(query1, organisation.OrganisationEmail, "organisations"); err != nil {
		return -1, err
	}

	if err := tx.QueryRow(query2, organisation.OrganisationMainAdminID, organisation.OrganisationName, organisation.OrganisationEmail, organisation.OrganisationPhoneNumber, organisation.OrganisationPassword).Scan(&orgID); err != nil {
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

func (q *Query) DeleteMainAdmin(mainAdminEmail string) (int, error) {
	query := "DELETE FROM main_admin WHERE main_admin_email = $1"
	if _, err := q.db.Exec(query, mainAdminEmail); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func (q *Query) DeleteOrganisation(organisationEmail string) (int, error) {
	query := "DELETE FROM organisations WHERE email = $1"
	if _, err := q.db.Exec(query, organisationEmail); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}
