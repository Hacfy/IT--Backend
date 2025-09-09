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

type DepartmentRepo struct {
	db *sql.DB
}

func NewDepartmentRepo(db *sql.DB) *DepartmentRepo {
	return &DepartmentRepo{
		db: db,
	}
}

func (dr *DepartmentRepo) CreateWorkspace(e echo.Context) (int, int, error) {
	status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return status, -1, err
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, -1, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, -1, fmt.Errorf("invalid user details")
	}

	var new_workspace models.CreateWorkspaceModel

	if err := e.Bind(&new_workspace); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, -1, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_workspace); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, -1, fmt.Errorf("failed to validate request")
	}

	status, workspace_id, err := query.CreateWorkspace(new_workspace, claims.UserID)
	if err != nil {
		log.Printf("error while storing Workspace data in DB: %v", err)
		return status, -1, fmt.Errorf("unable to create workspace at the moment, please try again later")
	}

	return status, workspace_id, nil
}

func (dr *DepartmentRepo) DeleteWorkspace(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var workspace models.DeleteWorkspaceModel

	if err := e.Bind(&workspace); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(workspace); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	workspace.WorkspaceName = strings.ToLower(workspace.WorkspaceName)

	status, err = query.DeleteWorkspace(workspace, claims.UserID)
	if err != nil {
		log.Printf("error while deleting the workspace %v: %v", workspace.WorkspaceID, err)
		return status, fmt.Errorf("error while deleting the workspace, please try again later")
	}

	return status, nil
}

func (dr *DepartmentRepo) RaiseIssue(e echo.Context) (int, int, error) {
	status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return status, -1, err
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, -1, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, -1, fmt.Errorf("invalid user details")
	}

	var Issue models.IssueModel

	if err := e.Bind(&Issue); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, -1, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(Issue); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, -1, fmt.Errorf("failed to validate request")
	}

	status, IssueID, err := query.RaiseIssue(Issue)
	if err != nil {
		log.Printf("error while storing Issue data in DB: %v", err)
		return status, -1, fmt.Errorf("unable to create issue at the moment, please try again later")
	}

	return status, IssueID, nil
}

func (dr *DepartmentRepo) RequestNewUnits(e echo.Context) (int, map[int]int, error) {
	Status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return Status, map[int]int{}, err
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, map[int]int{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, map[int]int{}, fmt.Errorf("invalid user details")
	}

	var new_unit models.RequestNewUnitModel

	if err := e.Bind(&new_unit); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, map[int]int{}, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_unit); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, map[int]int{}, fmt.Errorf("failed to validate request")
	}

	ok, err = query.CheckIfWarehouseIDExistsInTheDepartmentsBranch(new_unit.WarehouseID, claims.UserID)
	if err != nil {
		log.Printf("Error checking warehouse details: %v", err)
		return http.StatusInternalServerError, map[int]int{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid warehouse details")
		return http.StatusUnauthorized, map[int]int{}, fmt.Errorf("invalid warehouse details")
	}

	RequestIDs := make(map[int]int)

	for Key, Value := range new_unit.ComponentIDNoOfUnits {
		prefix, ok, err := query.CheckIfComponentIDExists(Key, new_unit.WarehouseID)
		if err != nil {
			log.Printf("Error checking component details: %v", err)
			return http.StatusInternalServerError, map[int]int{}, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid component details")
			return http.StatusUnauthorized, map[int]int{}, fmt.Errorf("invalid component details")
		}

		Status, requestID, err := query.RequestNewUnits(new_unit.DepartmentID, new_unit.WorkspaceID, new_unit.WarehouseID, Key, Value, prefix, claims.UserID)
		if err != nil {
			log.Printf("error while requesting new units: %v", err)
			return Status, map[int]int{}, fmt.Errorf("unable to request new units at the moment, please try again later")
		}

		RequestIDs[Key] = requestID
	}

	return http.StatusCreated, RequestIDs, nil
}

func (dr *DepartmentRepo) GetAllDepartmentRequests(e echo.Context) (int, []models.AllRequestsModel, error) {
	Status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return Status, []models.AllRequestsModel{}, fmt.Errorf("invalid use token")
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, []models.AllRequestsModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.AllRequestsModel{}, fmt.Errorf("invalid user details")
	}

	var getAllRequests models.GetAllRequestsModel

	err = e.Bind(&getAllRequests)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, []models.AllRequestsModel{}, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(getAllRequests); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, []models.AllRequestsModel{}, fmt.Errorf("failed to validate request")
	}

	requests, err := query.GetAllRequests(getAllRequests.DepartmentID)
	if err != nil {
		log.Printf("error while fetching requests: %v", err)
		return http.StatusInternalServerError, []models.AllRequestsModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, requests, nil
}

func (dr *DepartmentRepo) GetDepartmentRequestDetails(e echo.Context) (int, models.RequestDetailsModel, error) {
	Status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return Status, models.RequestDetailsModel{}, fmt.Errorf("invalid use token")
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, models.RequestDetailsModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, models.RequestDetailsModel{}, fmt.Errorf("invalid user details")
	}

	var GetRequestDetails models.GetRequestDetailsModel

	err = e.Bind(&GetRequestDetails)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, models.RequestDetailsModel{}, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(GetRequestDetails); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, models.RequestDetailsModel{}, fmt.Errorf("failed to validate request")
	}

	request, err := query.GetRequestDetails(GetRequestDetails)
	if err != nil {
		log.Printf("error while fetching request details: %v", err)
		return http.StatusInternalServerError, models.RequestDetailsModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, request, nil
}

func (dr *DepartmentRepo) DeleteIssue(e echo.Context) (int, error) {
	Status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return Status, err
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var deleteIssue models.DeleteIssueModel

	if err := e.Bind(&deleteIssue); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(deleteIssue); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	if deleteIssue.IssueID <= 0 {
		log.Printf("invalid issue id")
		return http.StatusBadRequest, fmt.Errorf("invalid issue id")
	}

	departmentID, err := query.GetDepartmentID(claims.UserID)
	if err != nil {
		log.Printf("error while getting department id: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if exists, err := query.CheckIfIssueIDExistsUnderDepartment(deleteIssue.IssueID, departmentID); err != nil {
		log.Printf("error while checking if issue exists: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !exists {
		log.Printf("issue with id %v does not exist", deleteIssue.IssueID)
		return http.StatusBadRequest, fmt.Errorf("issue with id %v does not exist", deleteIssue.IssueID)
	}

	status, err := query.DeleteIssue(deleteIssue.IssueID, claims.UserID)
	if err != nil {
		log.Printf("error while deleting issue: %v", err)
		return status, fmt.Errorf("database error")
	}

	return status, nil
}

func (dr *DepartmentRepo) DeleteRequest(e echo.Context) (int, error) {
	Status, claims, err := utils.VerifyUserToken(e, "department_heads", dr.db)
	if err != nil {
		return Status, err
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department_heads", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var deleteRequest models.DeleteRequestModel

	if err := e.Bind(&deleteRequest); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(deleteRequest); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	if deleteRequest.RequestID <= 0 {
		log.Printf("invalid request id")
		return http.StatusBadRequest, fmt.Errorf("invalid request id")
	}

	departmentID, err := query.GetDepartmentID(claims.UserID)
	if err != nil {
		log.Printf("error while getting department id: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if exists, err := query.CheckIfRequestIDExistsUnderDepartment(deleteRequest.RequestID, departmentID); err != nil {
		log.Printf("error while checking if request exists: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !exists {
		log.Printf("request with id %v does not exist", deleteRequest.RequestID)
		return http.StatusBadRequest, fmt.Errorf("request with id %v does not exist", deleteRequest.RequestID)
	}

	status, err := query.DeleteRequest(deleteRequest.RequestID, claims.UserID)
	if err != nil {
		log.Printf("error while deleting request: %v", err)
		return status, fmt.Errorf("database error")
	}

	return status, nil
}
