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
	status, claims, err := utils.VerifyUserToken(e, "department", dr.db)
	if err != nil {
		return status, -1, err
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details:", err)
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
	status, claims, err := utils.VerifyUserToken(e, "department", dr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details:", err)
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
	status, claims, err := utils.VerifyUserToken(e, "department", dr.db)
	if err != nil {
		return status, -1, err
	}
	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "department", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details:", err)
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
