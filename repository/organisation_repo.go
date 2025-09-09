package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
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
	status, claims, err := utils.VerifyUserToken(e, "organisations", or.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(or.db)

	ok, err := query.VerifyUser(claims.UserEmail, "organisations", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
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

	SuperAdmin.SuperAdminName = strings.ToLower(new_sa.SuperAdminName)

	SuperAdmin.Org_ID = claims.UserID

	SuperAdmin.SuperAdminPassword = hash

	SuperAdmin.SuperAdminID, err = query.CreateSuperAdmin(SuperAdmin)
	if err != nil {
		log.Printf("error while storing SuperAdmin data in DB: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to register SuperAdmin at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", SuperAdmin.SuperAdminEmail)
		if err := utils.SendLoginCredentials(SuperAdmin.SuperAdminEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", SuperAdmin.SuperAdminEmail, err)
		}
		log.Printf("credentials sent to %v", SuperAdmin.SuperAdminEmail)
	}()

	return http.StatusCreated, nil
}

func (or *OrgRepo) DeleteSuperAdmin(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "organisations", or.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(or.db)

	ok, err := query.VerifyUser(claims.UserEmail, "organisations", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var del_sa models.DeleteSuperAdminModel

	if err := e.Bind(&del_sa); err != nil {
		fmt.Printf("failed to decode request : %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(del_sa); err != nil {
		fmt.Printf("failed to validate request : %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	if status, err := query.DeleteSuperAdmin(del_sa.SuperAdminEmail); err != nil {
		log.Printf("error while deleting the user %v: %v", del_sa.SuperAdminEmail, err)
		return status, fmt.Errorf("error while deleting superAdmin, please try again later \n %v", err)
	}

	return http.StatusNoContent, nil

}

func (or *OrgRepo) GetAllSuperAdmins(e echo.Context) (int, []models.AllSuperAdminsDetailsModel, error) {
	status, claims, err := utils.VerifyUserToken(e, "organisations", or.db)
	if err != nil {
		return status, []models.AllSuperAdminsDetailsModel{}, err
	}

	query := database.NewDBinstance(or.db)

	ok, err := query.VerifyUser(claims.UserEmail, "organisations", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, []models.AllSuperAdminsDetailsModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.AllSuperAdminsDetailsModel{}, fmt.Errorf("invalid user details")
	}

	var getAllSuperAdmins models.GetAllSuperAdminsModel

	err = e.Bind(&getAllSuperAdmins)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, []models.AllSuperAdminsDetailsModel{}, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(getAllSuperAdmins); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, []models.AllSuperAdminsDetailsModel{}, fmt.Errorf("failed to validate request")
	}

	superAdmins, err := query.GetAllSuperAdmins(getAllSuperAdmins.OrganisationID)
	if err != nil {
		log.Printf("error while fetching superAdmins: %v", err)
		return http.StatusInternalServerError, []models.AllSuperAdminsDetailsModel{}, err
	}

	return http.StatusOK, superAdmins, nil

}

func (or *OrgRepo) ReassignSuperAdmin(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "organisations", or.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(or.db)

	ok, err := query.VerifyUser(claims.UserEmail, "organisations", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var reassignSuperAdmin models.ReassignSuperAdminModel

	err = e.Bind(&reassignSuperAdmin)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(reassignSuperAdmin); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	if err := query.CheckIfSuperAdminExists(reassignSuperAdmin.OldSuperAdminID, claims.UserID); err != nil {
		log.Printf("error while checking if superAdmin exists: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if err := query.CheckIfSuperAdminExists(reassignSuperAdmin.NewSuperAdminID, claims.UserID); err != nil {
		log.Printf("error while checking if superAdmin exists: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	status, err = query.ReassignSuperAdmin(reassignSuperAdmin, claims.UserID)
	if err != nil {
		log.Printf("error while reassigning superAdmin: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return status, nil
}
