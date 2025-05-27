package database

import (
	"database/sql"
	"fmt"
	"log"
)

type Query struct {
	db *sql.DB
}

func NewDBinstance(db *sql.DB) *Query {
	return &Query{
		db: db,
	}
}

func (db *Query) InitialiseDBqueries() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS main_admin (
			main_admin_id SERIAL PRIMARY KEY,
			main_admin_email VARCHAR(50) NOT NULL,
			main_admin_password VARCHAR(256) NOT NULL
		)`,
		`DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'userlevel') THEN
				CREATE TYPE userLevel AS ENUM (
					'organisations',
					'super_admins', 
					'branch_heads', 
					'department_heads'
					'workspace_heads'
				);
			END IF;
		END $$;`,
		`CREATE TABLE IF NOT EXISTS users (
			user_email VARCHAR(50) PRIMARY KEY,
			user_level userLevel NOT NULL,
			ever_logged_in BOOLEAN NOT NULL DEFAULT FALSE,
			latest_token TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS organisations (
			id SERIAL PRIMARY KEY,
			main_admin_id INTEGER NOT NULL,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) UNIQUE NOT NULL,
			phone_number VARCHAR(10) NOT NULL,
			password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_organisations_main_admin_id FOREIGN KEY (main_admin_id) REFERENCES main_admin(main_admin_id) ON UPDATE CASCADE,
			CONSTRAINT fk_organisations_email FOREIGN KEY (email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS super_admins (
			id SERIAL PRIMARY KEY,
			org_id INTEGER NOT NULL,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL,
			password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_super_admin_org_id FOREIGN KEY (org_id) REFERENCES organisations(id) ON DELETE CASCADE,
			CONSTRAINT fk_super_admin_super_admin_email FOREIGN KEY (email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS branches (
			branch_id SERIAL PRIMARY KEY,
			org_id INTEGER NOT NULL,
			super_admin_id INTEGER NOT NULL,
			branch_name VARCHAR(50) NOT NULL,
			branch_location VARCHAR(100) NOT NULL,
			CONSTRAINT fk_branch_org_id FOREIGN KEY (org_id) REFERENCES organisations(id) ON DELETE CASCADE,
			CONSTRAINT fk_branch_super_admin_id FOREIGN KEY (super_admin_id) REFERENCES super_admins(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS branch_heads (
			id SERIAL PRIMARY KEY,
			branch_id INTEGER UNIQUE NOT NULL,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL,
			password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_branch_heads_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON DELETE CASCADE,
			CONSTRAINT fk_branch_heads_email FOREIGN KEY (email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS departments (
			department_id SERIAL PRIMARY KEY,
			branch_id INTEGER NOT NULL,
			department_name VARCHAR(50) NOT NULL,
			CONSTRAINT fk_department_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS department_heads (
			id SERIAL PRIMARY KEY,
			department_id INTEGER UNIQUE NOT NULL,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL,
			password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_department_heads_department_id FOREIGN KEY (department_id) REFERENCES departments(department_id) ON DELETE CASCADE,
			CONSTRAINT fk_department_heads_email FOREIGN KEY (email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_super_admins (
			id SERIAL PRIMARY KEY,
			super_admin_id INTEGER NOT NULL,
			org_id INTEGER NOT NULL,
			email VARCHAR(50) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_organisations (
			id SERIAL PRIMARY KEY,
			org_id INTEGER NOT NULL,
			email VARCHAR(50) NOT NULL,
			main_admin_id INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_branches (
			id SERIAL PRIMARY KEY,
			branch_id INTEGER NOT NULL,
			super_admin_id INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_branch_heads (
			id SERIAL PRIMARY KEY,
			branch_id INTEGER NOT NULL,
			branch_head_id INTEGER NOT NULL,
			email VARCHAR(50) NOT NULL,
			deleted_by INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS warehouses (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL,
			password VARCHAR(256) NOT NULL,
			branch_id INTEGER NOT NULL,
			CONSTRAINT fk_warehouses_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON DELETE CASCADE, 
			CONSTRAINT fk_warehouses_email FOREIGN KEY (email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_department_heads (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			department_head_id INTEGER NOT NULL,
			email VARCHAR(50) NOT NULL,
			deleted_by INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_departments (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			branch_id INTEGER NOT NULL,
			deleted_by INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_warehouse_heads (
			id SERIAL PRIMARY KEY,
			warehouse_id INTEGER NOT NULL,
			email VARCHAR(50) NOT NULL,
			deleted_by INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS workespaces (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			workspace_name VARCHAR(50) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_workspaces (
			id SERIAL PRIMARY KEY,
			workspace_id INTEGER NOT NULL,
			department_id INTEGER NOT NULL,
			deleted_by INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS components (
			id SERIAL PRIMARY KEY,
			name VARCHAR(30) NOT NULL,
			prefix VARCHAR(3) NOT NULL UNIQUE,
			warehouse_id INTEGER NOT NULL,
			a_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT fk_component_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE
		)`,
	}

	tx, err := db.db.Begin()
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

	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			log.Printf("error while executing query %d: %v", i, err)
			return err
		}
	}

	return nil

}
