package database

import (
	"fmt"
	"log"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) CreateBranch(branch models.CreateBranchModel, superAdminID int, hashedPassword string) error {
	var branch_org_id, branch_id int
	query1 := "SELECT org_id FROM super_admins WHERE id = $1"
	query2 := "INSERT INTO branches(org_id, super_admin_id, branch_name, branch_location) VALUES($1, $2, $3, $4) RETURNING branch_id"
	query3 := "INSERT INTO users(user_email, user_level) VALUES($1, $2)"
	query4 := "INSERT INTO branch_heads(branch_id, name, email, password) VALUES($1, $2, $3, $4)"

	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("error while initialising DB: %v", err)
		return fmt.Errorf("something went wrong while processing your request. Please try again later")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	if err := tx.QueryRow(query1, superAdminID).Scan(&branch_org_id); err != nil {
		return err
	}

	if err := tx.QueryRow(query2, branch_org_id, superAdminID, branch.BranchName, branch.BranchLocation).Scan(&branch_id); err != nil {
		return err
	}

	if _, err := tx.Exec(query3, branch.BranchHeadEmail, "branch_heads"); err != nil {
		return err
	}

	if _, err := tx.Exec(query4, branch_id, branch.BranchHeadName, branch.BranchHeadEmail, hashedPassword); err != nil {
		return err
	}

	return nil

}
