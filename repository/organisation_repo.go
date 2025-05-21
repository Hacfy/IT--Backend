package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type OrgRepo struct {
	db *sql.DB
}

func NewOrgRepo(db *sql.DB) *OrgRepo {
	return &OrgRepo{
		db: db,
	}
}

func (or *OrgRepo) CreateSuperAdmin(e echo.Context) (int, error) {
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

	if claims.UserType != "organisations" {
		log.Printf("invalid userType, required userType %v given %v", "organisations", claims.UserType)
		return http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	}

	query := database.NewDBinstance(or.db)

	ok, err = query.VerifyUser(claims.UserEmail, "organisations", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details:", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var new_sa models.CreateSuperAdminModel

	if err := e.Bind(&new_sa); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_sa); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	password, err := utils.GeneratePassword()
	if err != nil {
		log.Printf("error while generating password: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("failed to generate password, please try again later")
	}

	for !utils.StrongPasswordValidator(password) {
		password, err = utils.GeneratePassword()
		if err != nil {
			log.Printf("error while generating password: %v", err)
			return http.StatusInternalServerError, fmt.Errorf("failed to generate password, please try again later")
		}
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("failed to secure your password, please try again")
	}

	var SuperAdmin models.SuperAdminModel

	SuperAdmin.Name = strings.ToLower(new_sa.SuperAdminName)

	SuperAdmin.Org_ID = claims.UserID

	SuperAdmin.Password = hash

	SuperAdmin.ID, err = query.CreateSuperAdmin(SuperAdmin)
	if err != nil {
		log.Printf("error while storing SuperAdmin data in DB: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to register SuperAdmin at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", SuperAdmin.Email)
		utils.SendLoginCredentials(SuperAdmin.Email, password)
		log.Printf("credentials sent to %v", SuperAdmin.Email)
	}()

	return http.StatusCreated, nil
}
