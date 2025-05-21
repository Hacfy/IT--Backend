package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/labstack/echo/v4"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) UserLogin(e echo.Context) (int, string, string, string, error) {
	var req_user models.LoginModel

	if err := e.Bind(&req_user); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, "", "", "", fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(req_user); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, "", "", "", fmt.Errorf("failded to validate request")
	}

	query := database.NewDBinstance(ar.db)

	userType, ok, err := query.GetUserType(req_user.Email)
	if err != nil {
		log.Printf("Error checking user details:", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid user credentials")
	}

	db_password, db_id, ok, err := query.GetUserPasswordID(req_user.Email, userType)
	if err != nil {
		log.Printf("Error checking user details:", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid user credentials")
	}

	if err := utils.CheckPassword(req_user.Password, db_password); err != nil {
		log.Printf("wrong password %v: %v", req_user.Email, err)
		return http.StatusBadRequest, "", "", "", fmt.Errorf("invalid user credentials")
	}

	accessToken, err := utils.GenerateCookieToken(req_user.Email, userType, db_id, time.Now().Local().Add(24*time.Hour).Unix())
	if err != nil {
		log.Printf("error while generating token for user %s: %v", req_user.Email, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	refreshToken, err := utils.GenerateCookieToken(req_user.Email, userType, db_id, time.Now().Local().Add(7*24*time.Hour).Unix())
	if err != nil {
		log.Printf("error while generating token for user %s: %v", req_user.Email, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	token, err := utils.GenerateUserToken(req_user.Email, userType, db_id, time.Now().Local().Add(7*24*time.Hour).Unix())
	if err != nil {
		log.Printf("error while generating token for user %s: %v", req_user.Email, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	return http.StatusOK, accessToken, refreshToken, token, nil

}
