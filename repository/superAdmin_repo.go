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

type SuperAdminRepo struct {
	db *sql.DB
}

func NewSuperAdminRepo(db *sql.DB) *SuperAdminRepo {
	return &SuperAdminRepo{
		db: db,
	}
}

func (sa *SuperAdminRepo) CreateBranch(e echo.Context) (int, error) {
	var tokenModel models.UserTokenModel

	tokenString := e.Request().Header.Get("Authorization")
	if tokenString == "" {
		log.Printf("missgin token")
		return http.StatusUnauthorized, fmt.Errorf("missing token")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &tokenModel, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("invalid token: %v", err)
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.UserTokenModel)
	if !(ok && token.Valid) {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	if claims.UserType != "super_admin" {
		log.Printf("invalid userType, required userType %v given %v", "super_admin", claims.UserType)
		return http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	}

	query := database.NewDBinstance(sa.db)

	ok, err = query.VerifyUser(claims.UserEmail, "super_admin", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var branch models.CreateBranchModel

	if err := e.Bind(&branch); err != nil {
		log.Printf("failed to decode request :%v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(branch); err != nil {
		log.Printf("failed to validate request :%v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	branch.BranchHeadName = strings.ToLower(branch.BranchHeadName)

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

	if err := query.CreateBranch(branch, claims.UserID, hash); err != nil {
		log.Printf("error while storing Branch Data in DB :%v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to create branch at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", branch.BranchHeadEmail)
		if err := utils.SendLoginCredentials(branch.BranchHeadEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", branch.BranchHeadEmail, err)
		}
		log.Printf("credentials sent to %v", branch.BranchHeadEmail)
	}()

	return http.StatusCreated, nil
}

func (sa *SuperAdminRepo) DeleteBranch(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "super_admin", sa.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(sa.db)

	ok, err := query.VerifyUser(claims.UserEmail, "super_admin", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var branch models.DeleteBranchModel

	if err := e.Bind(&branch); err != nil {
		log.Printf("failed to decode request :%v", err)
		return http.StatusBadRequest, fmt.Errorf("ivalid request format")
	}

	if err := validate.Struct(branch); err != nil {
		log.Printf("failed to validate request :%v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	branch.BrachName = strings.ToLower(branch.BrachName)

	if status, err := query.DeleteBranch(branch, claims.UserID); err != nil {
		log.Printf("error while deleting the branch %v: %v", branch.BranchID, err)
		return status, fmt.Errorf("error while deleting the branch")
	}

	return http.StatusNoContent, nil
}

func (sa *SuperAdminRepo) UpdateBranchHead(e echo.Context) (int, error) {
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
	if !(ok && token.Valid) {
		log.Printf("token expired or not of UserTokenModel")
		return http.StatusUnauthorized, fmt.Errorf("invalid token")
	}

	if claims.UserType != "super_admin" {
		log.Printf("invalid userType, required userType %v given %v", "super_admin", claims.UserType)
		return http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	}

	query := database.NewDBinstance(sa.db)

	ok, err = query.VerifyUser(claims.UserEmail, "super_admin", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var branchHead models.UpdateBranchHeadModel

	if err := e.Bind(&branchHead); err != nil {
		log.Printf("failed to decode request :%v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(branchHead); err != nil {
		log.Printf("failed to validate request :%v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	branchHead.NewBranchHeadName = strings.ToLower(branchHead.NewBranchHeadName)

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

	status, err := query.UpdateBranchHead(branchHead, claims.UserID, hash)
	if err != nil {
		log.Printf("error while deleting branchHead %v: %v", branchHead.BranchHeadID, err)
		return status, fmt.Errorf("unable to delete the branch head at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", branchHead.NewBranchHeadEmail)
		if err := utils.SendLoginCredentials(branchHead.NewBranchHeadEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", branchHead.NewBranchHeadEmail, err)
		}
		log.Printf("credentials sent to %v", branchHead.NewBranchHeadEmail)
	}()

	return http.StatusCreated, nil
}
