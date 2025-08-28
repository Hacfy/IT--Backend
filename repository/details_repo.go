package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/labstack/echo/v4"
)

type DetailsRepo struct {
	db *sql.DB
}

func NewDetailsRepo(db *sql.DB) *DetailsRepo {
	return &DetailsRepo{db: db}
}

// returns Departments, status, Total_Departments, Page, Limit, error
func (dr *DetailsRepo) GetAllDepartmentsRepo(e echo.Context) ([]models.AllDepartmentsModel, int, int, int, int, error) {
	var Sort models.SortModel
	Sort.Limit, _ = strconv.Atoi(e.QueryParam("limit"))
	if Sort.Limit <= 0 || Sort.Limit > 100 {
		Sort.Limit = 10
	}
	Sort.Page, _ = strconv.Atoi(e.QueryParam("page"))
	if Sort.Page <= 0 {
		Sort.Page = 1
	}
	Sort.Offset = (Sort.Page - 1) * Sort.Limit

	Sort.Order = e.QueryParam("order")
	if Sort.Order != "asc" && Sort.Order != "desc" {
		Sort.Order = "asc"
	}

	Sort.SortBy = e.QueryParam("sortBy")
	if Sort.SortBy == "" {
		Sort.SortBy = "department_id"
	}

	allowed := map[string]bool{"department_name": true, "department_id": true}
	if !allowed[Sort.SortBy] {
		Sort.SortBy = "department_id"
	}

	Sort.Search = e.QueryParam("search")
	role, ok := e.Get("userType").(string)
	if !ok {
		return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user credentials")
	}
	userID, ok := e.Get("userID").(int)
	if !ok {
		return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user credentials")
	}
	userEmail, ok := e.Get("userEmail").(string)
	if !ok {
		return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user credentials")
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(userEmail, role, userID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
	}

	var Request models.GetAllDepartmentsModel

	if err := e.Bind(&Request); err != nil {
		log.Printf("failed to decode request: %v", err)
		return []models.AllDepartmentsModel{}, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(Request); err != nil {
		log.Printf("failed to validate request: %v", err)
		return []models.AllDepartmentsModel{}, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("failed to validate request")
	}

	status, Departments, Total_Departments, err := query.GetAllDepartments(Request.BranchID, Sort)
	if err != nil {
		return []models.AllDepartmentsModel{}, status, Total_Departments, Sort.Page, Sort.Limit, err
	}

	return Departments, status, Total_Departments, Sort.Page, Sort.Limit, nil

}

func (dr *DetailsRepo) GetDepartmentIssues(e echo.Context) (int, []models.DepartmentIssuesModel, int, int, int, error) {
	var Sort models.SortModel
	Sort.Limit, _ = strconv.Atoi(e.QueryParam("limit"))
	if Sort.Limit <= 0 || Sort.Limit > 100 {
		Sort.Limit = 10
	}
	Sort.Page, _ = strconv.Atoi(e.QueryParam("page"))
	if Sort.Page <= 0 {
		Sort.Page = 1
	}
	Sort.Offset = (Sort.Page - 1) * Sort.Limit

	Sort.Order = e.QueryParam("order")
	if Sort.Order != "asc" && Sort.Order != "desc" {
		Sort.Order = "asc"
	}

	Sort.SortBy = e.QueryParam("sortBy")
	if Sort.SortBy == "" {
		Sort.SortBy = "created_at"
	}

	allowed := map[string]bool{"created_at": true}
	if !allowed[Sort.SortBy] {
		Sort.SortBy = "created_at"
	}

	Sort.Search = e.QueryParam("search")
	role, ok := e.Get("userType").(string)
	if !ok {
		return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}

	userID, ok := e.Get("userID").(int)
	if !ok {
		return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}

	userEmail, ok := e.Get("userEmail").(string)
	if !ok {
		return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(userEmail, role, userID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
	}

	var Request models.GetDepartmentIssuesModel

	if err := e.Bind(&Request); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(Request); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("failed to validate request")
	}

	status, Issues, Total_Issues, err := query.GetDepartmentIssues(Request.DepartmentID, Sort)
	if err != nil {
		return status, []models.DepartmentIssuesModel{}, Total_Issues, Sort.Page, Sort.Limit, err
	}

	return status, Issues, Total_Issues, Sort.Page, Sort.Limit, nil
}

func (dr *DetailsRepo) GetDepartmentWorkspaces(e echo.Context) ([]models.DepartmentWorkspaceModel, int, int, int, int, error) {
	var Sort models.SortModel
	Sort.Limit, _ = strconv.Atoi(e.QueryParam("limit"))
	if Sort.Limit <= 0 || Sort.Limit > 100 {
		Sort.Limit = 10
	}
	Sort.Page, _ = strconv.Atoi(e.QueryParam("page"))
	if Sort.Page <= 0 {
		Sort.Page = 1
	}
	Sort.Offset = (Sort.Page - 1) * Sort.Limit

	Sort.Order = e.QueryParam("order")
	if Sort.Order != "asc" && Sort.Order != "desc" {
		Sort.Order = "asc"
	}

	Sort.SortBy = e.QueryParam("sortBy")
	if Sort.SortBy == "" {
		Sort.SortBy = "department_id"
	}

	allowed := map[string]bool{"department_name": true, "department_id": true}
	if !allowed[Sort.SortBy] {
		Sort.SortBy = "department_id"
	}

	Sort.Search = e.QueryParam("search")
	role, ok := e.Get("userType").(string)
	if !ok {
		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}
	userID, ok := e.Get("userID").(int)
	if !ok {
		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}
	userEmail, ok := e.Get("userEmail").(string)
	if !ok {
		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(userEmail, role, userID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
	}

	var Request models.GetDepartmentWorkspacesModel

	if err := e.Bind(&Request); err != nil {
		log.Printf("failed to decode request: %v", err)
		return nil, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(Request); err != nil {
		log.Printf("failed to validate request: %v", err)
		return nil, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("failed to validate request")
	}

	status, Workspaces, Total_Workspaces, err := query.GetAllWorkspaces(Request.DepartmentID, Sort)
	if err != nil {
		return nil, status, -1, Sort.Page, Sort.Limit, err
	}

	return Workspaces, status, Total_Workspaces, Sort.Page, Sort.Limit, nil
}

// get all branches
func (dr *DetailsRepo) GetAllBranches(e echo.Context) ([]models.AllBranchesModel, int, int, int, int, error) {
	var Sort models.SortModel
	Sort.Limit, _ = strconv.Atoi(e.QueryParam("limit"))
	if Sort.Limit <= 0 || Sort.Limit > 100 {
		Sort.Limit = 10
	}
	Sort.Page, _ = strconv.Atoi(e.QueryParam("page"))
	if Sort.Page <= 0 {
		Sort.Page = 1
	}
	Sort.Offset = (Sort.Page - 1) * Sort.Limit

	Sort.Order = e.QueryParam("order")
	if Sort.Order != "asc" && Sort.Order != "desc" {
		Sort.Order = "asc"
	}

	Sort.SortBy = e.QueryParam("sortBy")
	if Sort.SortBy == "" {
		Sort.SortBy = "branch_id"
	}

	allowed := map[string]bool{"branch_name": true, "branch_id": true}
	if !allowed[Sort.SortBy] {
		Sort.SortBy = "branch_id"
	}

	Sort.Search = e.QueryParam("search")
	role, ok := e.Get("userType").(string)
	if !ok {
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}
	userID, ok := e.Get("userID").(int)
	if !ok {
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}
	userEmail, ok := e.Get("userEmail").(string)
	if !ok {
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(userEmail, role, userID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
	}

	var Request models.GetAllBranchesModel

	if err := e.Bind(&Request); err != nil {
		log.Printf("failed to decode request: %v", err)
		return []models.AllBranchesModel{}, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(Request); err != nil {
		log.Printf("failed to validate request: %v", err)
		return []models.AllBranchesModel{}, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("failed to validate request")
	}

	status, Branches, Total_Branches, err := query.GetAllBranches(Request.SuperAdminID, Sort)
	if err != nil {
		return []models.AllBranchesModel{}, status, -1, Sort.Page, Sort.Limit, err
	}

	return Branches, status, Total_Branches, Sort.Page, Sort.Limit, nil
}

// get branch details

// func (dr *DetailsRepo) GetBranchDetails(e echo.Context) ([]models.BranchDetailsModel, int, int, int, int, error) {
// 	var Sort models.SortModel
// 	Sort.Limit, _ = strconv.Atoi(e.QueryParam("limit"))
// 	if Sort.Limit <= 0 || Sort.Limit > 100 {
// 		Sort.Limit = 10
// 	}
// 	Sort.Page, _ = strconv.Atoi(e.QueryParam("page"))
// 	if Sort.Page <= 0 {
// 		Sort.Page = 1
// 	}
// 	Sort.Offset = (Sort.Page - 1) * Sort.Limit

// 	Sort.Order = e.QueryParam("order")
// 	if Sort.Order != "asc" && Sort.Order != "desc" {
// 		Sort.Order = "asc"
// 	}

// 	// Sort.SortBy = e.QueryParam("sortBy")
// 	// if Sort.SortBy == "" {
// 	// 	Sort.SortBy = "branch_id"
// 	// }
// 	Sort.Search = e.QueryParam("search")
// 	role, ok := e.Get("userType").(string)
// 	if !ok {
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
// 	}
// 	userID, ok := e.Get("userID").(int)
// 	if !ok {
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
// 	}
// 	userEmail, ok := e.Get("userEmail").(string)
// 	if !ok {
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
// 	}

// 	query := database.NewDBinstance(dr.db)

// 	ok, err := query.VerifyUser(userEmail, role, userID)
// 	if err != nil {
// 		log.Printf("Error checking user details:", err)
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
// 	} else if !ok {
// 		log.Printf("Invalid user details")
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
// 	}

// }

// get all warehouses

// func (dr *DetailsRepo) GetAllWarehouses(e echo.Context) ([]models.AllWarehousesModel, int, int, int, int, error) {
// 	var Sort models.SortModel
// 	Sort.Limit, _ = strconv.Atoi(e.QueryParam("limit"))
// 	if Sort.Limit <= 0 || Sort.Limit > 100 {
// 		Sort.Limit = 10
// 	}
// 	Sort.Page, _ = strconv.Atoi(e.QueryParam("page"))
// 	if Sort.Page <= 0 {
// 		Sort.Page = 1
// 	}
// 	Sort.Offset = (Sort.Page - 1) * Sort.Limit

// 	Sort.Order = e.QueryParam("order")
// 	if Sort.Order != "asc" && Sort.Order != "desc" {
// 		Sort.Order = "asc"
// 	}

// 	Sort.SortBy = e.QueryParam("sortBy")
// 	if Sort.SortBy == "" {
// 		Sort.SortBy = "warehouse_id"
// 	}

// 	allowed := map[string]bool{"warehouse_name": true, "warehouse_id": true}
// 	if !allowed[Sort.SortBy] {
// 		Sort.SortBy = "warehouse_id"
// 	}

// 	Sort.Search = e.QueryParam("search")
// 	role, ok := e.Get("userType").(string)
// 	if !ok {
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
// 	}
// 	userID, ok := e.Get("userID").(int)
// 	if !ok {
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
// 	}
// 	userEmail, ok := e.Get("userEmail").(string)
// 	if !ok {
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
// 	}

// 	query := database.NewDBinstance(dr.db)

// 	ok, err := query.VerifyUser(userEmail, role, userID)
// 	if err != nil {
// 		log.Printf("Error checking user details:", err)
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
// 	} else if !ok {
// 		log.Printf("Invalid user details")
// 		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
// 	}

// 	var Request models.GetAllWarehousesModel

// 	if err := e.Bind(&Request); err != nil {
// 		log.Printf("failed to decode request: %v", err)
// 		return nil, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
// 	}

// 	if err := validate.Struct(Request); err != nil {
// 		log.Printf("failed to validate request: %v", err)
// 		return nil, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("failed to validate request")
// 	}

// 	status, warehouses, total, err := query.GetAllWarehouses(Request.BranchID, Sort)
// 	if err != nil {
// 		return nil, status, -1, Sort.Page, Sort.Limit, err
// 	}

// 	return warehouses, status, total, Sort.Page, Sort.Limit, nil

// }

// // get workspace details
// // get department details
// // get warehouse details
