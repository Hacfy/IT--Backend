package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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

	if err = tx.QueryRow(query1, superAdminID).Scan(&branch_org_id); err != nil {
		return err
	}

	if err = tx.QueryRow(query2, branch_org_id, superAdminID, branch.BranchName, branch.BranchLocation).Scan(&branch_id); err != nil {
		return err
	}

	if _, err = tx.Exec(query3, branch.BranchHeadEmail, "branch_heads"); err != nil {
		return err
	}

	if _, err = tx.Exec(query4, branch_id, branch.BranchHeadName, branch.BranchHeadEmail, hashedPassword); err != nil {
		return err
	}

	return nil

}

func (q *Query) DeleteBranch(branch models.DeleteBranchModel, superAdminID int) (int, error) {
	query1 := "CALL delete_branch($1, $2)"

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

	if _, err := tx.Exec(query1, branch.BranchID, superAdminID); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil

}

func (q *Query) UpdateBranchHead(branchHead models.UpdateBranchHeadModel, superAdminID int, password string) (int, error) {
	query1 := "DELETE FROM branch_heads WHERE id =$1 RETURNING branch_id"
	query2 := "DELETE FROM users WHERE user_email = $1"
	query3 := "INSERT INTO deleted_branch_heads(branch_id, branch_head_id, email, deleted_by) VALUES($1, $2, $3, $4)"
	query4 := "INSERT INTO users(user_email, user_level) VALUES($1, $2)"
	query5 := "INSERT INTO branch_heads(branch_id, name, email, password) VALUES($1, $2, $3, $4)"

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

	var branch_id int

	if err = tx.QueryRow(query1, branchHead.BranchHeadID).Scan(&branch_id); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, err
	}

	if _, err = tx.Exec(query2, branchHead.BranchHeadEmail); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err = tx.Exec(query3, branch_id, branchHead.BranchHeadID, branchHead.BranchHeadEmail, superAdminID); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err = tx.Exec(query4, branchHead.NewBranchHeadEmail, "branch_heads"); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err = tx.Exec(query5, branch_id, branchHead.NewBranchHeadName, branchHead.NewBranchHeadEmail, password); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil

}
