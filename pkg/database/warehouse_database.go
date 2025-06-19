package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
)

func (q *Query) IfPrefixExists(prefix string) bool {
	var exists bool
	q.db.QueryRow("SELECT EXISTS(SELECT 1 FROM warehouses WHERE prefix = $1)", prefix).Scan(&exists)
	return exists
}

func (q *Query) IfComponentExists(name string, warehouse_id int) bool {
	var exists bool
	q.db.QueryRow("SELECT EXISTS(SELECT 1 FROM components WHERE name = $1 AND warehouse_id = $2)", name, warehouse_id).Scan(&exists)
	return exists
}

func (q *Query) CheckIfComponentIDExists(component_id, warehouse_id int) (string, bool, error) {
	query := "SELECT prefix FROM components WHERE id = $1 AND warehouse_id = $2"
	var prefix string
	err := q.db.QueryRow(query, component_id, warehouse_id).Scan(&prefix)
	if err != nil {
		return "", false, err
	}
	return prefix, true, nil
}

func (q *Query) CreateComponent(name, prefix string, warehouse_id int) (int, error) {
	var id int
	query1 := "INSERT INTO components(name, prefix, warehouse_id) VALUES($1, $2, $3) RETURNING id"
	query2 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s_units (
		id SERIAL PRIMARY KEY,
		component_id INTEGER NOT NULL,
		warehouse_id INTEGER NOT NULL,
		warranty_date TIMESTAMPTZ NOT NULL,
		status unit_status DEFAULT 'working',
		cost NUMERIC(10, 2) NOT NULL,
		maintainance_cost NUMERIC(10, 2) DEFAULT 0,
		CONSTRAINT fk_%s_units_component_id FOREIGN KEY (component_id) REFERENCES components(id) ON DELETE CASCADE,
		CONSTRAINT fk_%s_units_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE
	)`, prefix, prefix, prefix)
	query3 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s_units_assigned (
		id INTEGER PRIMARY KEY,
		department_id INTEGER NOT NULL,
		workspace_id INTEGER NOT NULL,
		CONSTRAINT fk_%s_units_department_id FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
		CONSTRAINT fk_%s_units_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		CONSTRAINT fk_%s_units_id FOREIGN KEY (id) REFERENCES %s_units(id) ON DELETE CASCADE
	)`, prefix, prefix, prefix, prefix)

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

	if err := tx.QueryRow(query1, name, prefix, warehouse_id).Scan(&id); err != nil {
		return -1, err
	}

	if _, err := tx.Exec(query2); err != nil {
		return -1, err
	}

	if _, err := tx.Exec(query3); err != nil {
		return -1, err
	}

	return id, nil
}

func (q *Query) DeleteComponent(del_component models.DeleteComponentModel, warehouse_id int) (int, error) {
	query1 := "DELETE FROM components WHERE id = $1 RETURNING id"
	query2 := "INSERT INTO deleted_components(component_id, warehouse_id, component_name, prefix, deleted_by) VALUES($1, $2, $3, $4)"

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

	if _, err := tx.Exec(query1, del_component.ComponentID); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		return http.StatusInternalServerError, err
	}

	if _, err := tx.Exec(query2, del_component.ComponentID, del_component.ComponentName, del_component.Prefix, warehouse_id); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

func (q *Query) CreateComponentUnit(warranty_date time.Time, cost float32, prifix string, warehouse_id, number, component_id int) (int, error) {
	query1 := fmt.Sprintf("INSERT INTO %s_units(component_id, warehouse_id, warranty_date, cost) VALUES($1, $2, $3, $4)", prifix)

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

	for _ = range number {
		_, err := tx.Exec(query1, component_id, warehouse_id, warranty_date, cost)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusCreated, nil

}

func (q *Query) AssignUnitWorkspace(workspace_id int, unit_id []int, prefix string) (int, error) {
	query1 := "SELECT department_id FROM workspaces WHERE id = $1"
	query2 := fmt.Sprintf("INSERT INTO %s_units_assigned(department_id, workspace_id, id) VALUES($1, $2, $3)", prefix)

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

	var department_id int

	if err := tx.QueryRow(query1, workspace_id).Scan(&department_id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no matching department found")
			return http.StatusNotFound, fmt.Errorf("no matching data found")
		}
		log.Printf("error while getting department id: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	for _, unit := range unit_id {
		if _, err := tx.Exec(query2, department_id, workspace_id, unit); err != nil {
			log.Printf("error while assigning units: %v", err)
			return http.StatusInternalServerError, fmt.Errorf("database error")
		}
	}

	return http.StatusOK, nil

}

// func (q *Query) GetAllIssues(warehouseID int)
