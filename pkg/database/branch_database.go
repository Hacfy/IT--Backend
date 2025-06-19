package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) CreateDepartment(department models.CreateDepartmentModel, branch_head_id int, password string) error {
	query0 := "SELECT branch_id FROM branch_heads WHERE id = $1"
	query1 := "INSERT INTO departments(branch_id, department_name) VALUES($1, $2) RETURNING department_id"
	query2 := "INSERT INTO users(user_email, user_level) VALUES($1, 2)"
	query3 := "INSERT INTO department_heads(department_id, name, email, password) VALUES($1, $2, $3, $4)"

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

	var department_id, branch_id int

	if err := tx.QueryRow(query0, branch_head_id).Scan(&branch_id); err != nil {
		return err
	}

	if err := tx.QueryRow(query1, branch_id, department.DepartmentName).Scan(&department_id); err != nil {
		return err
	}

	if _, err := tx.Exec(query2, department.DepartmentHeadEmail, "department_heads"); err != nil {
		return err
	}

	if _, err := tx.Exec(query3, department_id, department.DepartmentHeadName, department.DepartmentHeadEmail, password); err != nil {
		return err
	}

	return nil
}

func (q *Query) CreateWarehouse(warehouse models.CreateWarehouseModel, branchHeadID int, password string) (int, error) {
	query0 := "SELECT branch_id FROM branch_heads WHERE id = $1"
	query1 := "INSERT INTO users(user_email, user_level) VALUES($1, 2)"
	query2 := "INSERT INTO warehouses(name, email, branch_id, password) VALUES($1, $2, $3, $4)"
	var warehouse_id, branch_id int

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

	if err := tx.QueryRow(query0, branchHeadID).Scan(&branch_id); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}

	if err := tx.QueryRow(query1, warehouse.WarehouseUserEmail, "warehouses").Scan(&warehouse_id); err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = tx.Exec(query2, warehouse.WarehouseUserName, warehouse.WarehouseUserEmail, branch_id, password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}

func (q *Query) UpdateDepartmentHead(department_head models.UpdateDepartmentHeadModel, branch_head_id int, password string) (int, error) {
	query1 := "DELETE FROM department_heads WHERE email =$1 RETURNING department_id, id"
	query2 := "DELETE FROM users WHERE user_email = $1"
	query3 := "INSERT INTO deleted_department_heads(department_id, department_head_id, email, deleted_by) VALUES($1, $2, $3, $4)"
	query4 := "INSERT INTO users(user_email, user_level) VALUES($1, 2)"
	query5 := "INSERT INTO department_heads(department_id, name, email, password) VALUES($1, $2, $3, $4)"

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

	var department_id, department_head_id int

	if err := tx.QueryRow(query1, department_head.DepartmentHeadEmail).Scan(&department_id, &department_head_id); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query2, department_head.DepartmentHeadEmail); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query3, department_id, department_head_id, department_head.DepartmentHeadEmail, branch_head_id); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query4, department_head.NewDepartmentHeadEmail, "department_heads"); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query5, department_id, department_head.NewDepartmentHeadName, department_head.NewDepartmentHeadEmail, password); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (q *Query) UpdateWarehouseHead(warehouse_head models.UpdateWarehouseHeadModel, branch_head_id int, password string) (int, error) {
	query1 := "INSERT INTO users(user_email, user_level) VALUES($1, 2)"
	query2 := "UPDATE warehouses SET name = $1, email = $2, password = $3 WHERE id = $4"
	query3 := "INSERT INTO deleted_warehouse_heads(warehouse_id, email, deleted_by) VALUES($1, $2, $3)"
	query4 := "DELETE FROM users WHERE user_email = $1"

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

	var warehouse_id int

	if _, err := tx.Exec(query1, warehouse_head.NewWarehouseHeadEmail, "warehouses"); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := tx.QueryRow(query2, warehouse_head.NewWarehouseHeadName, warehouse_head.NewWarehouseHeadEmail, password, warehouse_head.WarehouseID).Scan(&warehouse_id); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query3, warehouse_id, warehouse_head.NewWarehouseHeadEmail, branch_head_id); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query4, warehouse_head.WarehouseHeadEmail); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
