package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
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
	var req_user models.UserLoginModel

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
		log.Printf("Error checking user type:", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user type")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid user credentials")
	}

	T, err := query.CheckUserLoggedIn(req_user.Email)
	if err != nil {
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("error while checking user details")
	}

	if !T {
		return http.StatusFound, "", "", "", fmt.Errorf("redirect")
	}

	db_password, db_name, db_id, ok, err := query.GetUserPasswordID(req_user.Email, userType)
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

	token_iat := time.Now().Local()

	if err := query.UpdateUserTokenTimestamp(req_user.Email, token_iat); err != nil {
		log.Printf("error while updating user token timestamp: %v", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("unable to update user token timestamp at the moment, please try again later")
	}

	token_unix := token_iat.Unix()

	accessToken, err := utils.GenerateCookieToken(req_user.Email, userType, db_id, time.Now().Local().Add(24*time.Hour).Unix(), token_unix)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", req_user.Email, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	refreshToken, err := utils.GenerateCookieToken(req_user.Email, userType, db_id, time.Now().Local().Add(7*24*time.Hour).Unix(), token_unix)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", req_user.Email, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	token, err := utils.GenerateUserToken(req_user.Email, userType, db_name, db_id, time.Now().Local().Add(7*24*time.Hour).Unix(), token_unix)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", req_user.Email, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	return http.StatusOK, accessToken, refreshToken, token, nil

}

func (ar *AuthRepo) ChangeUserPassword(e echo.Context) (int, string, string, string, error) {

	var tokenModel models.UserTokenModel

	tokenString := e.Request().Header.Get("Authorization")
	if tokenString == "" {
		log.Printf("missgin token")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("missing token")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &tokenModel, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("invalid token: %v", err)
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.UserTokenModel)
	if (ok && token.Valid) != true {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ar.db)

	Time, err := query.GetLatestTokenTime(claims.UserEmail, claims.UserType)
	if err != nil || Time.IsZero() {
		log.Printf("error while getting latest token time: %v", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("database error")
	}

	if Time.After(claims.IssuedAt.Time) {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid token")
	}

	var req_user models.ChangePasswordModel

	if err := e.Bind(&req_user); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, "", "", "", fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(req_user); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, "", "", "", fmt.Errorf("failded to validate request")
	}

	if !utils.StrongPasswordValidator(req_user.NewPassword) {
		log.Printf("invalid password")
		return http.StatusBadRequest, "", "", "", fmt.Errorf("invalid password")
	}

	hash, err := utils.HashPassword(req_user.NewPassword)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("failed to secure your password, please try again")
	}

	status, err := query.ChangeUserPassword(hash, claims.UserEmail, claims.UserType)
	if err != nil {
		log.Printf("error while storing new password in DB: %v", err)
		return status, "", "", "", fmt.Errorf("unable to update password at the moment, please try again later")
	}

	token_iat := time.Now().Local()

	if err := query.UpdateUserTokenTimestamp(claims.UserEmail, token_iat); err != nil {
		log.Printf("error while updating user token timestamp: %v", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("unable to update user token timestamp at the moment, please try again later")
	}

	token_unix := token_iat.Unix()

	accessToken, err := utils.GenerateCookieToken(claims.UserEmail, claims.UserType, claims.UserID, time.Now().Local().Add(24*time.Hour).Unix(), token_unix)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", claims.UserEmail, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	refreshToken, err := utils.GenerateCookieToken(claims.UserEmail, claims.UserType, claims.UserID, time.Now().Local().Add(7*24*time.Hour).Unix(), token_unix)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", claims.UserEmail, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	Token, err := utils.GenerateUserToken(claims.UserEmail, claims.UserType, claims.UserName, claims.UserID, time.Now().Local().Add(7*24*time.Hour).Unix(), token_unix)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", claims.UserEmail, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	return http.StatusOK, accessToken, refreshToken, Token, nil
}

func (ar *AuthRepo) UserLogout(e echo.Context) (int, error) {

	var tokenModel models.UserTokenModel

	tokenString := e.Request().Header.Get("Authorization")
	if tokenString == "" {
		log.Printf("missgin token")
		return http.StatusUnauthorized, fmt.Errorf("missing token")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &tokenModel, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("invalid token: %v", err)
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.UserTokenModel)
	if (ok && token.Valid) != true {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ar.db)

	Time, err := query.GetLatestTokenTime(claims.UserEmail, claims.UserType)
	if err != nil || Time.IsZero() {
		log.Printf("error while getting latest token time: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if Time.After(claims.IssuedAt.Time) {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	var req_user models.UserLogoutModel

	if err := e.Bind(&req_user); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(req_user); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, fmt.Errorf("failded to validate request")
	}

	if claims.UserEmail != req_user.Email {
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	userType, ok, err := query.GetUserType(req_user.Email)
	if err != nil {
		if !ok {
			return http.StatusUnauthorized, fmt.Errorf("invalid user details")
		}
		log.Printf("Error checking user details:", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if claims.UserType != userType {
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	token_iat := time.Now().Local()

	if err := query.UpdateUserTokenTimestamp(req_user.Email, token_iat); err != nil {
		log.Printf("error while updating user token timestamp: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to update user token timestamp at the moment, please try again later")
	}

	return http.StatusOK, nil

}
