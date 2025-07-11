package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func GenerateMainAdminToken(mainAdmin models.MainAdminModel) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"main_admin_id":    mainAdmin.MainAdminID,
		"main_admin_email": mainAdmin.MainAdminEmail,
		"exp":              time.Now().Local().Add(7 * 24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateOrganisationToken(organisation models.OrganisationModel) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"organisation_id":            organisation.OrganisationID,
		"organisation_email":         organisation.OrganisationEmail,
		"organisation_main_admin_id": organisation.OrganisationMainAdminID,
		"exp":                        time.Now().Local().Add(7 * 24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateCookieToken(userEmail, userType string, userID int, exp, iat int64) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"user_email": userEmail,
		"user_type":  userType,
		"iat":        iat,
		"exp":        exp,
		"iss":        "IT_INVENTORY",
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateUserToken(userEmail, userType, userName string, userID int, exp, iat int64) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"user_name":  userName,
		"user_email": userEmail,
		"user_type":  userType,
		"iat":        iat,
		"exp":        exp,
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyUserToken(e echo.Context, userType string, db *sql.DB) (int, models.UserTokenModel, error) {
	var tokenModel models.UserTokenModel

	tokenString := e.Request().Header.Get("Authorization")
	if tokenString == "" {
		log.Printf("missgin token")
		return http.StatusUnauthorized, tokenModel, fmt.Errorf("missing token")
	}

	token, err := jwt.ParseWithClaims(tokenString, &tokenModel, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("invalid token: %v", err)
		return http.StatusUnauthorized, tokenModel, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.UserTokenModel)
	if (ok && token.Valid) != true {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, tokenModel, fmt.Errorf("invalid token")
	}

	if claims.UserType != userType {
		log.Printf("invalid userType, required userType %v given %v", userType, claims.UserType)
		return http.StatusUnauthorized, tokenModel, fmt.Errorf("invalid credentials")
	}

	query := database.NewDBinstance(db)

	time, err := query.GetLatestTokenTime(claims.UserEmail, claims.UserType)
	if err != nil || time.IsZero() {
		log.Printf("error while getting latest token time: %v", err)
		return http.StatusInternalServerError, tokenModel, fmt.Errorf("database error")
	}

	if time.After(claims.IssuedAt.Time) {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, tokenModel, fmt.Errorf("invalid token")
	}

	return http.StatusOK, tokenModel, nil
}

func GenerateComponentToken(id int, name, prefix string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"component_id":   id,
		"component_name": name,
		"prefix":         prefix,
		"exp":            time.Now().Local().Add(7 * 24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("error while generating token: %v", err)
		return "", err
	}

	return signedToken, nil
}

func ParseToken(tokenStr string) (*models.UserTokenModel, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	Claims, ok := token.Claims.(*models.UserTokenModel)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return Claims, nil
}
