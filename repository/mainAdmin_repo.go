package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
			log.Fatalf("error while sending credentials to %v: %v", main_admin.MainAdminEmail, err)
		}
		log.Printf("credentials sent to %v", main_admin.MainAdminEmail)
	}()

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
		log.Printf("Error checking main admin details: %v", err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid main admin details")
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid main admin credentials")
	}

	if err = utils.CheckPassword(login_ma.MainAdminPassword, db_ma.MainAdminPassword); err != nil {
		log.Printf("wrong password %v: %v", login_ma.MainAdminEmail, err)
		return http.StatusUnauthorized, "", "", "", fmt.Errorf("invalid main admin password")
	}

	DB_iat := time.Now().Local().Unix()

	accessToken, err := utils.GenerateCookieToken(db_ma.MainAdminEmail, "main_admin", db_ma.MainAdminID, time.Now().Local().Add(24*time.Hour).Unix(), DB_iat)
	if err != nil {
		log.Printf("error while generating token for user %s: %v", db_ma.MainAdminEmail, err)
		return http.StatusInternalServerError, "", "", "", fmt.Errorf("error while generating token for user %s: %v", db_ma.MainAdminEmail, err)
	}

	refreshToken, err := utils.GenerateCookieToken(db_ma.MainAdminEmail, "main_admin", db_ma.MainAdminID, time.Now().Local().Add(7*24*time.Hour).Unix(), DB_iat)
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
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("invalid token: %v", err)
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.MainAdminTokenModel)
	if !(ok && token.Valid) {
		log.Printf("token expired or not of MainAdminTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ma.db)

	ok, err = query.VerifyMainAdmin(claims.MainAdminEmail, claims.MainAdminID)
	if err != nil {
		log.Printf("Error checking main admin details: %v", err)
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

	createOrganisationPassword := os.Getenv("ORGANISATION_PASSWORD")

	if new_org.CreateOrganisationPassword != createOrganisationPassword {
		log.Printf("wrong create_organisation_password")
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

	var organisation models.OrganisationModel

	organisation.OrganisationEmail = new_org.OrganisationEmail

	organisation.OrganisationName = strings.ToLower(new_org.OrganisationName)

	organisation.OrganisationPhoneNumber = new_org.OrganisationPhoneNumber

	organisation.OrganisationPassword = hash

	organisation.OrganisationMainAdminID = claims.MainAdminID

	organisation.OrganisationID, err = query.CreateOrganisation(organisation)
	if err != nil {
		log.Printf("error while storing organisation data in DB: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to register organisation at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", organisation.OrganisationEmail)
		if err := utils.SendLoginCredentials(organisation.OrganisationEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", organisation.OrganisationEmail, err)
		}
		log.Printf("credentials sent to %v", organisation.OrganisationEmail)
	}()

	return http.StatusCreated, nil
}

func (ma *MainAdminRepo) DeleteMainAdmin(e echo.Context) (int, error) {

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
	if !(ok && token.Valid) {
		log.Printf("token expired or not of MainAdminTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ma.db)

	ok, err = query.VerifyMainAdmin(claims.MainAdminEmail, claims.MainAdminID)
	if err != nil {
		log.Printf("Error checking main admin details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid main admin details")
		return http.StatusUnauthorized, fmt.Errorf("invalid main admin details")
	}

	var main_admin models.DeleteMainAdminModel

	if err = e.Bind(&main_admin); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(main_admin); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, fmt.Errorf("failded to validate request")
	}

	companyPassword := os.Getenv("COMPANY_PASSWORD")

	if main_admin.CompanyPassword != companyPassword {
		log.Printf("wrong company_password")
		return http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	}

	status, err := query.DeleteMainAdmin(main_admin.MainAdminEmail, main_admin.MainAdminID, tokenModel.MainAdminID)
	if err != nil {
		log.Printf("error while deleting the main admin from database %v: %v", main_admin.MainAdminEmail, err)
		return status, fmt.Errorf("error while deleting main admin")
	}

	return status, nil
}

func (ma *MainAdminRepo) DeleteOrganisation(e echo.Context) (int, error) {

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
	if !(ok && token.Valid) {
		log.Printf("token expired or not of MainAdminTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ma.db)

	ok, err = query.VerifyMainAdmin(claims.MainAdminEmail, claims.MainAdminID)
	if err != nil {
		log.Printf("Error checking main admin details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid main admin details")
		return http.StatusUnauthorized, fmt.Errorf("invalid main admin details")
	}

	var del_org models.DeleteOrganisationModel

	if err = e.Bind(&del_org); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(del_org); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, fmt.Errorf("failded to validate request")
	}

	deleteOrganisationPassword := os.Getenv("ORGANISATION_PASSWORD")

	if del_org.DeleteOrganisationPassword != deleteOrganisationPassword {
		log.Printf("wrong delete_organisation_password")
		return http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	}

	status, err := query.DeleteOrganisation(del_org.OrganisationEmail, del_org.OrganisationID, claims.MainAdminID)
	if err != nil {
		log.Printf("error while deleting the organisation %v: %v", del_org.OrganisationEmail, err)
		return status, fmt.Errorf("error while deleting organisation, please try again later")
	}

	return status, nil
}

func (ma *MainAdminRepo) GetAllOrganisations(e echo.Context) (int, []models.GetAllOrganisationsModel, error) {
	var tokenModel models.MainAdminTokenModel

	tokenString := e.Request().Header.Get("Authorization")
	if tokenString == "" {
		log.Printf("missgin token")
		return http.StatusUnauthorized, []models.GetAllOrganisationsModel{}, fmt.Errorf("missing token")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &tokenModel, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("invalid token: %v", err)
		return http.StatusUnauthorized, []models.GetAllOrganisationsModel{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.MainAdminTokenModel)
	if !(ok && token.Valid) {
		log.Printf("token expired or not of MainAdminTokenModel")
		return http.StatusUnauthorized, []models.GetAllOrganisationsModel{}, fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ma.db)

	ok, err = query.VerifyMainAdmin(claims.MainAdminEmail, claims.MainAdminID)
	if err != nil {
		log.Printf("Error checking main admin details: %v", err)
		return http.StatusInternalServerError, []models.GetAllOrganisationsModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid main admin details")
		return http.StatusUnauthorized, []models.GetAllOrganisationsModel{}, fmt.Errorf("invalid main admin details")
	}

	orgs, err := query.GetAllOrganisations(claims.MainAdminID)
	if err != nil {
		log.Println("Error while getting organisations:", err)
		return http.StatusInternalServerError, []models.GetAllOrganisationsModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, orgs, nil

}

func (ma *MainAdminRepo) GetAllMainAdmins(e echo.Context) (int, []models.AllMainAdminModel, error) {
	var tokenModel models.MainAdminTokenModel

	tokenString := e.Request().Header.Get("Authorization")
	if tokenString == "" {
		log.Printf("missgin token")
		return http.StatusUnauthorized, []models.AllMainAdminModel{}, fmt.Errorf("missing token")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &tokenModel, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("invalid token: %v", err)
		return http.StatusUnauthorized, []models.AllMainAdminModel{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.MainAdminTokenModel)
	if !(ok && token.Valid) {
		log.Printf("token expired or not of MainAdminTokenModel")
		return http.StatusUnauthorized, []models.AllMainAdminModel{}, fmt.Errorf("invalid token")
	}

	query := database.NewDBinstance(ma.db)

	ok, err = query.VerifyMainAdmin(claims.MainAdminEmail, claims.MainAdminID)
	if err != nil {
		log.Printf("Error checking main admin details: %v", err)
		return http.StatusInternalServerError, []models.AllMainAdminModel{}, fmt.Errorf("database error")
	}

	if !ok {
		log.Printf("Invalid main admin details")
		return http.StatusUnauthorized, []models.AllMainAdminModel{}, fmt.Errorf("invalid main admin details")
	}

	var deleteRequest models.GetAllMainAdminModel

	err = e.Bind(&deleteRequest)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, []models.AllMainAdminModel{}, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(deleteRequest); err != nil {
		log.Printf("failed to validate request %v", err)
		return http.StatusBadRequest, []models.AllMainAdminModel{}, fmt.Errorf("failded to validate request")
	}

	companyPassword := os.Getenv("COMPANY_PASSWORD")

	if deleteRequest.CompanyPassword != companyPassword {
		log.Printf("wrong company_password")
		return http.StatusUnauthorized, []models.AllMainAdminModel{}, fmt.Errorf("invalid credentials")
	}

	main_admins, err := query.GetAllMainAdmins()
	if err != nil {
		log.Println("Error while getting main admins:", err)
		return http.StatusInternalServerError, []models.AllMainAdminModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, main_admins, nil
}
