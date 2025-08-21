package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/labstack/echo/v4"
)

type BranchRepo struct {
	db *sql.DB
}

func NewBranchRepo(db *sql.DB) *BranchRepo {
	return &BranchRepo{
		db: db,
	}
}

func (br *BranchRepo) CreateDepartment(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "branch_heads", br.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(br.db)

	ok, err := query.VerifyUser(claims.UserEmail, "branch_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var new_department models.CreateDepartmentModel

	if err := e.Bind(&new_department); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_department); err != nil {
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

	if err := query.CreateDepartment(new_department, claims.UserID, hash); err != nil {
		log.Printf("error while storing Department data in DB: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("unable to create department at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", new_department.DepartmentHeadEmail)
		if err := utils.SendLoginCredentials(new_department.DepartmentHeadEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", new_department.DepartmentHeadEmail, err)
		}
		log.Printf("credentials sent to %v", new_department.DepartmentHeadEmail)
	}()

	return http.StatusCreated, nil
}

func (br *BranchRepo) CreateWarehouse(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "branch_heads", br.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(br.db)

	ok, err := query.VerifyUser(claims.UserEmail, "branch_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var new_warehouse models.CreateWarehouseModel

	if err := e.Bind(&new_warehouse); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_warehouse); err != nil {
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

	status, err = query.CreateWarehouse(new_warehouse, claims.UserID, hash)
	if err != nil {
		log.Printf("error while storing Warehouse data in DB: %v", err)
		return status, fmt.Errorf("unable to create warehouse at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", new_warehouse.WarehouseUserEmail)
		if err := utils.SendLoginCredentials(new_warehouse.WarehouseUserEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", new_warehouse.WarehouseUserEmail, err)
		}
		log.Printf("credentials sent to %v", new_warehouse.WarehouseUserEmail)
	}()

	return http.StatusCreated, nil
}

func (br *BranchRepo) UpdateDepartmentHead(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "branch_heads", br.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(br.db)

	ok, err := query.VerifyUser(claims.UserEmail, "branch_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var new_department_head models.UpdateDepartmentHeadModel

	if err := e.Bind(&new_department_head); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_department_head); err != nil {
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

	if status, err := query.UpdateDepartmentHead(new_department_head, claims.UserID, hash); err != nil {
		log.Printf("error while storing DepartmentHead data in DB: %v", err)
		return status, fmt.Errorf("unable to update department head at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", new_department_head.NewDepartmentHeadEmail)
		if err := utils.SendLoginCredentials(new_department_head.NewDepartmentHeadEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", new_department_head.NewDepartmentHeadEmail, err)
		}
		log.Printf("credentials sent to %v", new_department_head.NewDepartmentHeadEmail)
	}()

	return http.StatusOK, nil
}

func (br *BranchRepo) UpdateWarehouseHead(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "branch_heads", br.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(br.db)

	ok, err := query.VerifyUser(claims.UserEmail, "branch_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var new_warehouse_head models.UpdateWarehouseHeadModel

	if err := e.Bind(&new_warehouse_head); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_warehouse_head); err != nil {
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

	if status, err := query.UpdateWarehouseHead(new_warehouse_head, claims.UserID, hash); err != nil {
		log.Printf("error while storing WarehouseHead data in DB: %v", err)
		return status, fmt.Errorf("unable to update warehouse head at the moment, please try again later")
	}

	go func() {
		log.Printf("sending login credentials to %v", new_warehouse_head.NewWarehouseHeadEmail)
		if err := utils.SendLoginCredentials(new_warehouse_head.NewWarehouseHeadEmail, password); err != nil {
			log.Printf("error while sending login credentials to %v: %v", new_warehouse_head.NewWarehouseHeadEmail, err)
		}
		log.Printf("credentials sent to %v", new_warehouse_head.NewWarehouseHeadEmail)
	}()

	return http.StatusOK, nil
}

func (br *BranchRepo) DeleteDepartment(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "branch_heads", br.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(br.db)

	ok, err := query.VerifyUser(claims.UserEmail, "branch_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var department models.DeleteDepartmentModel

	if err := e.Bind(&department); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(department); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	ok, err = query.IsDepartmentUnderBranchHead(department.DepartmentID, claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	status, err = query.DeleteDepartment(department.DepartmentID, claims.UserID)
	if err != nil {
		log.Printf("error while deleting the department %v: %v", department.DepartmentID, err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return status, nil
}

func (br *BranchRepo) DeleteWarehouse(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "branch_heads", br.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(br.db)

	ok, err := query.VerifyUser(claims.UserEmail, "branch_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var warehouse models.DeleteWarehouseModel

	if err := e.Bind(&warehouse); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(warehouse); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	ok, err = query.VerifyUser(claims.UserEmail, "branch_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	UserPassword, _, _, ok, err := query.GetUserPasswordID(claims.UserEmail, claims.UserType)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	if err := utils.CheckPassword(warehouse.BranchHeadPassword, UserPassword); err != nil {
		log.Printf("wrong password %v: %v", warehouse.BranchHeadPassword, err)
		return http.StatusBadRequest, fmt.Errorf("invalid user details")
	}
	ok, err = query.IsWarehouseUnderBranchHead(warehouse.WarehouseID, claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	status, err = query.DeleteWarehouse(warehouse.WarehouseID, claims.UserID)
	if err != nil {
		log.Printf("error while deleting the warehouse %v: %v", warehouse.WarehouseID, err)
		return status, fmt.Errorf("database error")
	}

	return status, nil
}
