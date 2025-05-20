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
				);
			END IF;
		END $$;`,
		`CREATE TABLE IF NOT EXISTS users (
			user_email VARCHAR(50) PRIMARY KEY,
			user_level userLevel NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS organisations (
			org_id SERIAL PRIMARY KEY,
			main_admin_id INTEGER NOT NULL,
			org_name VARCHAR(50) NOT NULL,
			org_email VARCHAR(50) NOT NULL,
			org_password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_organisations_main_admin_id FOREIGN KEY (main_admin_id) REFERENCES main_admin(main_admin_id) ON UPDATE CASCADE,
			CONSTRAINT fk_organisations_email FOREIGN KEY (org_email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS super_admins (
			super_admin_id SERIAL PRIMARY KEY,
			org_id INTEGER NOT NULL,
			super_admin_email VARCHAR(50) NOT NULL,
			super_admin_password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_super_admin_org_id FOREIGN KEY (org_id) REFERENCES organisations(org_id) ON DELETE CASCADE,
			CONSTRAINT fk_super_admin_super_admin_email FOREIGN KEY (super_admin_email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS branches (
			branch_id SERIAL PRIMARY KEY,
			org_id INTEGER NOT NULL,
			super_admin_id INTEGER NOT NULL,
			branch_name VARCHAR(50) NOT NULL,
			branch_location VARCHAR(50) NOT NULL,
			CONSTRAINT fk_branch_org_id FOREIGN KEY (org_id) REFERENCES organisations(org_id) ON DELETE CASCADE,
			CONSTRAINT fk_branch_super_admin_id FOREIGN KEY (super_admin_id) REFERENCES super_admins(super_admin_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS branch_heads (
			branch_head_id SERIAL PRIMARY KEY,
			branch_id INTEGER NOT NULL,
			branch_head_name VARCHAR(50) NOT NULL,
			branch_head_email VARCHAR(50) NOT NULL,
			branch_head_password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_branch_heads_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON DELETE CASCADE,
			CONSTRAINT fk_branch_heads_email FOREIGN KEY (branch_head_email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS departments (
			department_id SERIAL PRIMARY KEY,
			branch_id INTEGER NOT NULL,
			department_name VARCHAR(50) NOT NULL,
			CONSTRAINT fk_department_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS department_heads (
			department_head_id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			department_head_name VARCHAR(50) NOT NULL,
			department_head_email VARCHAR(50) NOT NULL,
			department_head_password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_department_heads_department_id FOREIGN KEY (department_id) REFERENCES departments(department_id) ON DELETE CASCADE,
			CONSTRAINT fk_department_heads_email FOREIGN KEY (department_head_email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
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
