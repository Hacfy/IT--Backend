package utils

import (
	"os"
	"time"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/golang-jwt/jwt/v5"
)

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

func GenerateCookieToken(userEmail, userType string, userID int, exp int64) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"user_email": userEmail,
		"user_type":  userType,
		"iat":        time.Now(),
		"exp":        exp,
		"iss":        "IT_INVENTORY",
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateUserToken(userEmail, userType string, userID int, exp int64) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"user_email": userEmail,
		"user_type":  userType,
		"exp":        exp,
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
