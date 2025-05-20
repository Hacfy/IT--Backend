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
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

type MainAdminRepo struct {
	db *sql.DB
}

func NewMainAdminRepo(db *sql.DB) *MainAdminRepo {
	return &MainAdminRepo{
		db: db,
	}
}

func (ma *MainAdminRepo) CreateMainAdmin(e echo.Context) (int, error) {
	var create_ma models.CreateMainAdminModel

	var main_admin models.MainAdminModel

	query := database.NewDBinstance(ma.db)

	if err := e.Bind(&create_ma); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(create_ma); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, fmt.Errorf("failded to validate request")
	}

	companyPassword := os.Getenv("COMPANY_PASSWORD")

	log.Println(companyPassword)

	if create_ma.CompanyPassword != companyPassword {
		log.Printf("wrong company_password")
		return http.StatusUnauthorized, fmt.Errorf("invalid credentials")
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

	main_admin.MainAdminEmail = create_ma.MainAdminEmail

	main_admin.MainAdminPassword = hash

	main_admin.MainAdminID, err = query.CreateMainAdmin(main_admin)
	if err != nil {
		log.Printf("error while storing main_admin data in DB: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to register main_admin at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", main_admin.MainAdminEmail)
		if err := utils.SendLoginCredentials(main_admin.MainAdminEmail, password); err != nil {
			log.Fatal("error while sending credentials to %v: %v", main_admin.MainAdminEmail, err)
		}
		log.Printf("credentials sent to %v", main_admin.MainAdminEmail)
	}()

	// token, err := utils.GenerateMainAdminToken(main_admin)
	// if err != nil {
	// 	log.Printf("error while generating token for user %s: %v", main_admin.MainAdminEmail, err)
	// 	return "", fmt.Errorf("unable to generate token, please try again later")
	// }

	return http.StatusCreated, nil
}

func (ma *MainAdminRepo) LoginMainAdmin(e echo.Context) (int, string, string, string, error) {
	var login_ma models.LoginMainAdminModel

	if err := e.Bind(&login_ma); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, "", "", "", fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(login_ma); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, "", "", "", fmt.Errorf("failded to validate request")
	}

	query := database.NewDBinstance(ma.db)

	db_ma, ok, err := query.GetMainAdminCredentials(login_ma.MainAdminEmail)
	if err != nil {
		log.Printf("Error checking main admin details:", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid main admin details")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid main admin credentials")
	}

	if err = utils.CheckPassword(login_ma.MainAdminPassword, db_ma.MainAdminPassword); err != nil {
		log.Printf("wrong password %v: %v", login_ma.MainAdminEmail, err)
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid main admin credentials")
	}

	accessToken, err := utils.GenerateCookieToken(db_ma.MainAdminEmail, "main_admin", db_ma.MainAdminID, time.Now().Local().Add(24*time.Hour).Unix())
	if err != nil {
		log.Printf("error while generating token for user %s: %v", db_ma.MainAdminEmail, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	refreshToken, err := utils.GenerateCookieToken(db_ma.MainAdminEmail, "main_admin", db_ma.MainAdminID, time.Now().Local().Add(7*24*time.Hour).Unix())
	if err != nil {
		log.Printf("error while generating token for user %s: %v", db_ma.MainAdminEmail, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	token, err := utils.GenerateMainAdminToken(db_ma)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", db_ma.MainAdminEmail, err)
		return http.StatusInternalServerError, "", "", "", err
	}

	return http.StatusOK, accessToken, refreshToken, token, nil

}

func (ma *MainAdminRepo) CreateOrganisation(e echo.Context) (int, error) {

	var tokenModel models.MainAdminTokenModel

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

	claims, ok := token.Claims.(*models.MainAdminTokenModel)
	if (ok && token.Valid) != true {
		log.Printf("token expired or not of MainAdminTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ma.db)

	ok, err = query.VerifyMainAdmin(claims.MainAdminEmail, claims.MainAdminID)
	if err != nil {
		log.Printf("Error checking main admin details:", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid main admin details")
		return http.StatusUnauthorized, fmt.Errorf("invalid main admin details")
	}

	var new_org models.CreateOrganisationModel

	if err = e.Bind(&new_org); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_org); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, fmt.Errorf("failded to validate request")
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

	var organisation models.OrganisationModel

	organisation.OrganisationEmail = new_org.OrganisationEmail

	organisation.OrganisationName = new_org.OrganisationName

	organisation.OrganisationPassword = hash

	organisation.OrganisationMainAdminID = claims.MainAdminID

	organisation.OrganisationID, err = query.CreateOrganisation(organisation)
	if err != nil {
		log.Printf("error while storing organisation data in DB: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to register organisation at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", organisation.OrganisationEmail)
		utils.SendLoginCredentials(organisation.OrganisationEmail, password)
		log.Printf("credentials sent to %v", organisation.OrganisationEmail)
	}()

	// OrgToken, err := utils.GenerateOrganisationToken(organisation)
	// if err != nil {
	// 	log.Printf("error while generating token for organisaion %s: %v", organisation.OrganisationEmail, err)
	// 	return fmt.Errorf("unable to generate token, please try again later")
	// }

	return http.StatusCreated, nil
}
