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

	if _, err := tx.Exec(query1, superAdmin.Email, "super_admins"); err != nil {
		return -1, err
	}

	if err := tx.QueryRow(query2, superAdmin.Org_ID, superAdmin.Name, superAdmin.Email, superAdmin.Password).Scan(&sa_id); err != nil {
		return -1, err
	}

	return sa_id, nil

}

func (q *Query) DeleteSuperAdmin(superAdminEmail string)
