package database

import "database/sql"

func (q *Query) GetUserType(userEmail string) (string, bool, error) {
	var userType string
	query := "SELECT user_level FROM users WHERE user_email = $1"

	if err := q.db.QueryRow(query, userEmail).Scan(&userType); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return userType, true, nil
}

func (q *Query) GetUserPassword(userEmail, userType string) (string, int, bool, error) {
	var db_password string
	var db_id int

	query := "SELECT password, id FROM $1 WHERE email = $2"
	if err := q.db.QueryRow(query, userType, userEmail).Scan(&db_password, &db_id); err != nil {
		if err == sql.ErrNoRows {
			return "", -1, false, nil
		}
		return "", -1, false, err
	}
	return db_password, db_id, true, nil
}
