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
		`CREATE TABLE IF NOT EXISTS deleted_main_admins (
			id SERIAL PRIMARY KEY,
			main_admin_id INTEGER NOT NULL,
			main_admin_email VARCHAR(50) NOT NULL,
			deleted_by INTEGER NOT NULL
		)`,
		`DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'userlevel') THEN
				CREATE TYPE userLevel AS ENUM (
					'organisations',
					'super_admins', 
					'branch_heads', 
					'department_heads',
					'warehouses'
				);
			END IF;
		END $$;`,
		`DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'unit_status') THEN
				CREATE TYPE unit_status AS ENUM (
					'working',
					'repair',
					'not_working',
					'exit'
				);
			END IF;
		END $$;`,
		`DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'issue_status') THEN
				CREATE TYPE issue_status AS ENUM (
					'raised',
					'accepted',
					'resolved'
				);
			END IF;
		END $$;`,
		`DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'request_status') THEN
				CREATE TYPE request_status AS ENUM (
					'raised',
					'accepted',
					'declined'
				);
			END IF;
		END $$;`,
		`CREATE TABLE IF NOT EXISTS users (
			user_email VARCHAR(50) PRIMARY KEY,
			user_level userLevel NOT NULL,
			ever_logged_in BOOLEAN NOT NULL DEFAULT FALSE,
			latest_token TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_users (
			id SERIAL PRIMARY KEY,
			user_email VARCHAR(50) NOT NULL,
			user_level userLevel NOT NULL,
			ever_logged_in BOOLEAN,
			latest_token TIMESTAMPTZ,
			created_at TIMESTAMPTZ,
			deleted_by INTEGER NOT NULL,
			deleted_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS otps (
			email VARCHAR(50),
			otp VARCHAR(6) NOT NULL,
			time TIMESTAMPTZ NOT NULL
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
			email VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_super_admin_org_id FOREIGN KEY (org_id) REFERENCES organisations(id) ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT fk_super_admin_super_admin_email FOREIGN KEY (email) REFERENCES users(user_email) ON DELETE CASCADE ON UPDATE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS branches (
			branch_id SERIAL PRIMARY KEY,
			org_id INTEGER NOT NULL,
			super_admin_id INTEGER NOT NULL,
			branch_name VARCHAR(50) NOT NULL,
			branch_location VARCHAR(500) NOT NULL,
			CONSTRAINT fk_branch_org_id FOREIGN KEY (org_id) REFERENCES organisations(id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_branch_super_admin_id FOREIGN KEY (super_admin_id) REFERENCES super_admins(id) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS branch_heads (
			id SERIAL PRIMARY KEY,
			branch_id INTEGER UNIQUE NOT NULL,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL,
			password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_branch_heads_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_branch_heads_email FOREIGN KEY (email) REFERENCES users(user_email) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS departments (
			department_id SERIAL PRIMARY KEY,
			branch_id INTEGER NOT NULL,
			department_name VARCHAR(50) NOT NULL,
			CONSTRAINT fk_department_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS department_heads (
			id SERIAL PRIMARY KEY,
			department_id INTEGER UNIQUE NOT NULL,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL,
			password VARCHAR(256) NOT NULL,
			CONSTRAINT fk_department_heads_department_id FOREIGN KEY (department_id) REFERENCES departments(department_id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_department_heads_email FOREIGN KEY (email) REFERENCES users(user_email) ON UPDATE CASCADE ON DELETE CASCADE
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
			CONSTRAINT fk_warehouses_branch_id FOREIGN KEY (branch_id) REFERENCES branches(branch_id) ON UPDATE CASCADE ON DELETE CASCADE, 
			CONSTRAINT fk_warehouses_email FOREIGN KEY (email) REFERENCES users(user_email) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_department_heads (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			department_head_id INTEGER NOT NULL,
			email VARCHAR(50) NOT NULL,
			deleted_by INTEGER NOT NULL,
			deleted_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_departments (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			branch_id INTEGER NOT NULL,
			deleted_by INTEGER NOT NULL, 
			deleted_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_warehouse_heads (
			id SERIAL PRIMARY KEY,
			warehouse_id INTEGER NOT NULL,
			email VARCHAR(50) NOT NULL,
			deleted_by INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS workspaces (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			workspace_name VARCHAR(50) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_workpaces (
			id SERIAL PRIMARY KEY,
			workspace_id INTEGER NOT NULL,
			department_id INTEGER NOT NULL,
			deleted_by INTEGER NOT NULL,
			deleted_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS components (
			id SERIAL PRIMARY KEY,
			name VARCHAR(30) NOT NULL,
			prefix VARCHAR(3) NOT NULL UNIQUE,
			warehouse_id INTEGER NOT NULL,
			a_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT fk_component_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_components (
			id SERIAL PRIMARY KEY,
			component_id INTEGER NOT NULL,
			component_name VARCHAR(30) NOT NULL,
			prefix VARCHAR(3) NOT NULL,
			deleted_by INTEGER NOT NULL,
			deleted_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS issues (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			warehouse_id INTEGER NOT NULL,
			workspace_id INTEGER NOT NULL,
			unit_id INTEGER NOT NULL,
			unit_prefix VARCHAR(3) NOT NULL,
			issue VARCHAR(100) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			status issue_status DEFAULT 'raised',
			CONSTRAINT fk_issues_department_id FOREIGN KEY (department_id) REFERENCES departments(department_id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_issues_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_issues_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS resolved_issues (
			id SERIAL PRIMARY KEY,
			issue_id INTEGER NOT NULL,
			solution VARCHAR(250) NOT NULL,
			cost NUMERIC(10, 2) NOT NULL,
			resolved_by INTEGER NOT NULL,
			resolved_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT fk_resolved_issues_issue_id FOREIGN KEY (issue_id) REFERENCES issues(id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_resolved_issues_resolved_by FOREIGN KEY (resolved_by) REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_issues (
			id SERIAL PRIMARY KEY,
			issue_id INTEGER NOT NULL,
			department_id INTEGER NOT NULL,
			workspace_id INTEGER NOT NULL,
			unit_id INTEGER NOT NULL,
			unit_prefix VARCHAR(3) NOT NULL,
			issue VARCHAR(100) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			status issue_status DEFAULT 'raised',
			deleted_by INTEGER NOT NULL,
			deleted_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS deleted_units_assigned (
			id SERIAL PRIMARY KEY,
			unit_id INTEGER NOT NULL,
			department_id INTEGER NOT NULL,
			workspace_id INTEGER NOT NULL,
			assigned_at TIMESTAMPTZ,
			deleted_by INTEGER NOT NULL,
			deleted_at TIMESTAMPTZ DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS deleted_units (
			id SERIAL PRIMARY KEY,
			unit_id INTEGER NOT NULL,
			unit_prefix VARCHAR(3) NOT NULL,
			component_id INTEGER NOT NULL,
			warehouse_id INTEGER NOT NULL,
			deleted_by INTEGER NOT NULL,
			deleted_at TIMESTAMPTZ DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS  requests (
			id SERIAL PRIMARY KEY,
			department_id INTEGER NOT NULL,
			workspace_id INTEGER NOT NULL,
			warehouse_id INTEGER NOT NULL,
			component_id INTEGER NOT NULL,
			number_of_units INTEGER NOT NULL,
			prefix VARCHAR(3) NOT NULL,
			created_by INTEGER NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			status request_status DEFAULT 'raised',
			CONSTRAINT fk_requests_department_id FOREIGN KEY (department_id) REFERENCES departments(department_id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_requests_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_requests_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_requests_component_id FOREIGN KEY (component_id) REFERENCES components(id) ON UPDATE CASCADE ON DELETE CASCADE
		)`,
	}

	// CREATE EXTENSION IF NOT EXISTS pg_cron;

	// SELECT cron.schedule(
	// 	'cleanup_deleted_units',
	// 	'0 0 * * *',
	// 	$$
	// 	DELETE FROM deleted_units WHERE deleted_at < NOW() - INTERVAL '30 days';
	// 	$$
	// );

	queries = append(queries, "CREATE OR REPLACE PROCEDURE delete_department(dep_id INTEGER, deleter_id INTEGER) LANGUAGE plpgsql AS $$ DECLARE r_workspace RECORD; r_component RECORD; r_dep_head RECORD; has_assignments BOOLEAN; dyn_sql TEXT; BEGIN INSERT INTO deleted_departments(department_id, branch_id, deleted_by) SELECT department_id, branch_id, deleter_id FROM departments WHERE department_id = dep_id; INSERT INTO deleted_department_heads(department_id, department_head_id, email, deleted_by) SELECT department_id, id, email, deleter_id FROM department_heads WHERE department_id = dep_id; FOR r_workspace IN SELECT * FROM workspaces WHERE department_id = dep_id LOOP INSERT INTO deleted_workpaces(workspace_id, department_id, deleted_by) VALUES (r_workspace.id, r_workspace.department_id, deleter_id); FOR r_component IN SELECT DISTINCT c.id, c.name, c.prefix FROM components c WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = lower(c.prefix || '_units_assigned')) LOOP dyn_sql := format('SELECT EXISTS (SELECT 1 FROM %I_units_assigned WHERE department_id = $1 AND workspace_id = $2)', r_component.prefix); EXECUTE dyn_sql INTO has_assignments USING dep_id, r_workspace.id; IF has_assignments THEN INSERT INTO deleted_components(component_id, component_name, prefix, deleted_by) VALUES (r_component.id, r_component.name, r_component.prefix, deleter_id); EXECUTE format('INSERT INTO deleted_units_assigned(unit_id, department_id, workspace_id, assigned_at, deleted_by) SELECT id, department_id, workspace_id, assigned_at, $1 FROM %I_units_assigned WHERE department_id = $2 AND workspace_id = $3', r_component.prefix) USING deleter_id, dep_id, r_workspace.id; EXECUTE format('UPDATE %I_units SET status = ''not_assigned'' WHERE id IN (SELECT unit_id FROM deleted_units_assigned WHERE department_id = $1 AND workspace_id = $2 AND deleted_by = $3)', r_component.prefix) USING dep_id, r_workspace.id, deleter_id; END IF; END LOOP; END LOOP; DELETE FROM workspaces WHERE department_id = dep_id; FOR r_dep_head IN SELECT email FROM department_heads WHERE department_id = dep_id LOOP INSERT INTO deleted_users(user_email, user_level, ever_logged_in, latest_token, created_at, deleted_by) SELECT u.user_email, u.user_level, u.ever_logged_in, u.latest_token, u.created_at, deleter_id FROM users u WHERE u.user_email = r_dep_head.email; DELETE FROM users WHERE user_email = r_dep_head.email; END LOOP; DELETE FROM department_heads WHERE department_id = dep_id; DELETE FROM departments WHERE department_id = dep_id; END; $$;",
		"CREATE OR REPLACE PROCEDURE delete_component(comp_id INTEGER, deleter_id INTEGER) LANGUAGE plpgsql AS $$ DECLARE r_comp_details RECORD; units_table_name TEXT; units_assigned_table_name TEXT; BEGIN SELECT name, prefix INTO r_comp_details FROM components WHERE id = comp_id; IF NOT FOUND THEN RAISE EXCEPTION 'Component with ID % not found.', comp_id; END IF; units_table_name := lower(r_comp_details.prefix || '_units'); units_assigned_table_name := lower(r_comp_details.prefix || '_units_assigned'); INSERT INTO deleted_components(component_id, component_name, prefix, deleted_by) VALUES (comp_id, r_comp_details.name, r_comp_details.prefix, deleter_id); IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = units_assigned_table_name) THEN EXECUTE format('INSERT INTO deleted_units_assigned(unit_id, department_id, workspace_id, assigned_at, deleted_by) SELECT sua.id, sua.department_id, sua.workspace_id, sua.assigned_at, $1 FROM %I sua JOIN %I su ON sua.id = su.id WHERE su.component_id = $2', units_assigned_table_name, units_table_name) USING deleter_id, comp_id; END IF; IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = units_table_name) THEN EXECUTE format('INSERT INTO deleted_units(unit_id, unit_prefix, component_id, warehouse_id, deleted_by) SELECT id, $1, component_id, warehouse_id, $2 FROM %I WHERE component_id = $3', units_table_name) USING r_comp_details.prefix, deleter_id, comp_id; END IF; IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = units_table_name) THEN EXECUTE format('DELETE FROM %I WHERE component_id = $1', units_table_name) USING comp_id; END IF; DELETE FROM components WHERE id = comp_id; END $$;",
		"CREATE OR REPLACE PROCEDURE delete_warehouse(wh_id INTEGER, deleter_id INTEGER) LANGUAGE plpgsql AS $$ DECLARE r_warehouse_head_details RECORD; r_component RECORD; BEGIN SELECT id, email INTO r_warehouse_head_details FROM warehouses WHERE id = wh_id; IF NOT FOUND THEN RAISE EXCEPTION 'Warehouse head with ID % not found.', wh_id; END IF; INSERT INTO deleted_warehouse_heads(warehouse_id, email, deleted_by) VALUES (r_warehouse_head_details.id, r_warehouse_head_details.email, deleter_id); INSERT INTO deleted_users(user_email, user_level, ever_logged_in, latest_token, created_at, deleted_by) SELECT u.user_email, u.user_level, u.ever_logged_in, u.latest_token, u.created_at, deleter_id FROM users u WHERE u.user_email = r_warehouse_head_details.email; DELETE FROM users WHERE user_email = r_warehouse_head_details.email; FOR r_component IN SELECT id FROM components WHERE warehouse_id = wh_id LOOP CALL delete_component(r_component.id, deleter_id); END LOOP; DELETE FROM warehouses WHERE id = wh_id; END $$;",
		"CREATE OR REPLACE PROCEDURE delete_branch(br_id INTEGER, deleter_id INTEGER) LANGUAGE plpgsql AS $$ DECLARE r_branch_head RECORD; r_department RECORD; r_warehouse RECORD; BEGIN INSERT INTO deleted_branches(branch_id, super_admin_id, deleted_by) SELECT branch_id, super_admin_id, deleter_id FROM branches WHERE branch_id = br_id; FOR r_branch_head IN SELECT id, email FROM branch_heads WHERE branch_id = br_id LOOP INSERT INTO deleted_branch_heads(branch_id, branch_head_id, email, deleted_by) VALUES (br_id, r_branch_head.id, r_branch_head.email, deleter_id); INSERT INTO deleted_users(user_email, user_level, ever_logged_in, latest_token, created_at, deleted_by) SELECT u.user_email, u.user_level, u.ever_logged_in, u.latest_token, u.created_at, deleter_id FROM users u WHERE u.user_email = r_branch_head.email; DELETE FROM users WHERE user_email = r_branch_head.email; END LOOP; FOR r_department IN SELECT department_id FROM departments WHERE branch_id = br_id LOOP CALL delete_department(r_department.department_id, deleter_id); END LOOP; FOR r_warehouse IN SELECT id FROM warehouses WHERE branch_id = br_id LOOP CALL delete_warehouse(r_warehouse.id, deleter_id); END LOOP; DELETE FROM branch_heads WHERE branch_id = br_id; DELETE FROM departments WHERE branch_id = br_id; DELETE FROM warehouses WHERE branch_id = br_id; DELETE FROM branches WHERE branch_id = br_id; END $$;",
	)

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
