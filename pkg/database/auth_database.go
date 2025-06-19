package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (q *Query) CheckUserLoggedIn(userEmail string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE user_email = $1 AND ever_logged_in = true)"
	if err := q.db.QueryRow(query, userEmail).Scan(&exists); err != nil {
		log.Println("Query error:", err)
		return false, err
	}
	if !exists {
		_, err := q.db.Exec("UPDATE users SET ever_logged_in = true WHERE user_email = $1 AND ever_logged_in = false", userEmail)
		if err != nil {
			log.Println("Update error:", err)
			return false, err
		}
		return false, nil
	}
	return true, nil
}

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

func (q *Query) UpdateUserTokenTimestamp(email string, T time.Time) error {
	query := "UPDATE users SET latest_token = $1 WHERE user_email = $2"
	if _, err := q.db.Exec(query, T, email); err != nil {
		return err
	}
	return nil
}

func (q *Query) GetUserPasswordID(userEmail, userType string) (string, string, int, bool, error) {
	var db_password, db_name string
	var db_id int

	query := fmt.Sprintf("SELECT password, id, name FROM %s WHERE email = $1", userType)
	if err := q.db.QueryRow(query, userEmail).Scan(&db_password, &db_id, &db_name); err != nil {
		if err == sql.ErrNoRows {
			return "", "", -1, false, nil
		}
		return "", "", -1, false, err
	}
	return db_password, db_name, db_id, true, nil
}

func (q *Query) VerifyUser(userEmail, userType string, userID int) (bool, error) {
	var exists int
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE email = $1 AND id = $2", userType)
	if err := q.db.QueryRow(query, userEmail, userID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (q *Query) GetLatestTokenTime(userEmail, userType string) (time.Time, error) {
	var latestToken time.Time
	query := fmt.Sprintf("SELECT latest_token FROM users WHERE user_email = $1 AND user_level = $2", userEmail, userType)
	if err := q.db.QueryRow(query, userEmail, userType).Scan(&latestToken); err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return latestToken, nil
}

func (q *Query) ChangeUserPassword(newPassword, userEmail, userType string) (int, error) {
	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE email = $2", userType)
	if _, err := q.db.Exec(query, newPassword, userEmail); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
