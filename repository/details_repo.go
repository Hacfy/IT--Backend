package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/labstack/echo/v4"
)

type DetailsRepo struct {
	db *sql.DB
}

func NewDetailsRepo(db *sql.DB) *DetailsRepo {
	return &DetailsRepo{db: db}
}

// remove the details route and implement them in the specified repo

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

	BranchID, err := strconv.Atoi(e.QueryParam("branch_id"))
	if err != nil {
		log.Printf("error while parsing branch id: %v", err)
		return []models.AllDepartmentsModel{}, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
	}

	if BranchID == 0 {
		log.Printf("invalid branch id")
		return []models.AllDepartmentsModel{}, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid branch id")
	}

	switch role {
	case "branch_head":
		ok, err := query.CheckBranchHead(userID, BranchID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "super_admin":
		ok, err := query.CheckIfBranchUnderSuperAdmin(BranchID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "organization":
		ok, err := query.CheckIfBranchUnderorganizationAdmin(BranchID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	default:
		log.Printf("Invalid user role")
		return []models.AllDepartmentsModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user role")
	}

	status, Departments, Total_Departments, err := query.GetAllDepartments(BranchID, Sort)
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

	DepartmentID, err := strconv.Atoi(e.QueryParam("department_id"))
	if err != nil {
		log.Printf("error while parsing department id: %v", err)
		return http.StatusBadRequest, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
	}

	switch role {
	case "department_head":
		ok, err := query.CheckDepartmentHead(userID, DepartmentID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "branch_head":
		ok, err := query.CheckIfDepartmentUnderBranchHead(DepartmentID, userID)
		if err != nil {
			log.Printf("Error checking user details:: %v", err)
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "super_admin":
		ok, err := query.CheckIfDepartmentUnderSuperAdmin(DepartmentID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "organization_admin":
		ok, err := query.CheckIfDepartmentUnderorganizationAdmin(DepartmentID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	default:
		log.Printf("Invalid user role")
		return http.StatusUnauthorized, []models.DepartmentIssuesModel{}, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user role")
	}

	status, Issues, Total_Issues, err := query.GetDepartmentIssues(DepartmentID, Sort)
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

	DepartmentID, err := strconv.Atoi(e.QueryParam("department_id"))
	if err != nil {
		log.Printf("error while parsing department id: %v", err)
		return nil, http.StatusBadRequest, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid request format")
	}

	switch role {
	case "department_head":
		ok, err := query.CheckDepartmentHead(userID, DepartmentID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "branch_head":
		ok, err := query.CheckIfDepartmentUnderBranchHead(DepartmentID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "super_admin":
		ok, err := query.CheckIfDepartmentUnderSuperAdmin(DepartmentID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	case "organization_admin":
		ok, err := query.CheckIfDepartmentUnderorganizationAdmin(DepartmentID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
		}
	default:
		log.Printf("Invalid user role")
		return nil, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user role")
	}

	status, Workspaces, Total_Workspaces, err := query.GetAllWorkspaces(DepartmentID, Sort)
	if err != nil {
		return nil, status, -1, Sort.Page, Sort.Limit, err
	}

	return Workspaces, status, Total_Workspaces, Sort.Page, Sort.Limit, nil
}

// get all branches
func (dr *DetailsRepo) GetAllBranchesUnderSuperAdmin(e echo.Context) ([]models.AllBranchesModel, int, int, int, int, error) {
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
	userID, ok := e.Get("userID").(int)
	if !ok {
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}
	userEmail, ok := e.Get("userEmail").(string)
	if !ok {
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid use credentials")
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(userEmail, "super_admin", userID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return []models.AllBranchesModel{}, http.StatusUnauthorized, -1, Sort.Page, Sort.Limit, fmt.Errorf("invalid user details")
	}

	status, Branches, Total_Branches, err := query.GetAllBranches(userID, Sort)
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

func (dr *DetailsRepo) GetAllWarehouses(e echo.Context) ([]models.AllWarehousesModel, int, error) {

	role, ok := e.Get("userType").(string)
	if !ok {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid use credentials")
	}
	userID, ok := e.Get("userID").(int)
	if !ok {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid use credentials")
	}
	userEmail, ok := e.Get("userEmail").(string)
	if !ok {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid use credentials")
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(userEmail, role, userID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return nil, http.StatusUnauthorized, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	BranchID, err := strconv.Atoi(e.QueryParam("branch_id"))
	if err != nil {
		log.Printf("error while parsing branch id: %v", err)
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	switch role {
	case "branch_head":
		ok, err := query.CheckBranchHead(userID, BranchID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return nil, http.StatusUnauthorized, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return nil, http.StatusUnauthorized, fmt.Errorf("invalid user details")
		}
	case "super_admin":
		ok, err := query.CheckIfBranchUnderSuperAdmin(BranchID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return nil, http.StatusUnauthorized, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return nil, http.StatusUnauthorized, fmt.Errorf("invalid user details")
		}
	case "organization_admin":
		ok, err := query.CheckIfBranchUnderorganizationAdmin(BranchID, userID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return nil, http.StatusUnauthorized, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return nil, http.StatusUnauthorized, fmt.Errorf("invalid user details")
		}
	default:
		log.Printf("Invalid user role")
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid user role")
	}

	status, warehouses, err := query.GetAllWarehouses(BranchID)
	if err != nil {
		return nil, status, err
	}

	return warehouses, status, nil
}

// // get workspace details
// // get department details
// // get warehouse details

func (dr *DetailsRepo) GetAllOutOfWarentyUnitsInDepartment(e echo.Context) (int, []models.AllOutOfWarentyUnitsModel, int, int, int, error) {
	role, ok := e.Get("userType").(string)
	if !ok {
		return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid use credentials")
	}

	userEmail, ok := e.Get("userEmail").(string)
	if !ok {
		return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid use credentials")
	}

	userID, ok := e.Get("userID").(int)
	if !ok {
		return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid use credentials")
	}

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

	ComponentID, err := strconv.Atoi(e.QueryParam("component_id"))
	if err != nil {
		log.Printf("error while parsing component id: %v", err)
		return http.StatusBadRequest, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid component id")
	}

	DepartmentID, err := strconv.Atoi(e.QueryParam("department_id"))
	if err != nil {
		log.Printf("error while parsing department id: %v", err)
		return http.StatusBadRequest, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid department id")
	}

	query := database.NewDBinstance(dr.db)

	ok, err = query.VerifyUser(userEmail, role, userID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid user details")
	}

	switch role {
	case "department_head":
		ok, err := query.CheckDepartmentHead(userID, DepartmentID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid user details")
		}
	case "warehouses":
		ok, err := query.CheckWarehouseHead(userID, ComponentID)
		if err != nil {
			log.Printf("Error checking user details: %v", err)
			return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("database error")
		} else if !ok {
			log.Printf("Invalid user details")
			return http.StatusUnauthorized, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid user details")
		}
	}

	_, prefix, err := query.GetComponentNameAndPrefix(ComponentID)
	if err != nil {
		log.Printf("error while fetching assigned units: %v", err)
		return http.StatusInternalServerError, []models.AllOutOfWarentyUnitsModel{}, -1, -1, -1, fmt.Errorf("database error")
	}

	status, warehouses, total, err := query.GetAllOutOfWarentyUnitsInDepartment(DepartmentID, ComponentID, prefix, Sort.Limit, Sort.Offset)
	if err != nil {
		return status, []models.AllOutOfWarentyUnitsModel{}, total, -1, -1, err
	}

	return status, warehouses, total, Sort.Limit, Sort.Page, nil
}

func (dr *DetailsRepo) GetAllOutOfWarentyUnitsInWarehouse(e echo.Context) (int, []models.AllOutOfWarentyWarehouseModel, int, int, int, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", dr.db)
	if err != nil {
		return status, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, err
	}

	query := database.NewDBinstance(dr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, fmt.Errorf("invalid user details")
	}

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

	ComponentID, err := strconv.Atoi(e.QueryParam("component_id"))
	if err != nil {
		log.Printf("error while parsing component id: %v", err)
		return http.StatusBadRequest, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, fmt.Errorf("invalid component id")
	}

	if ComponentID <= 0 {
		log.Printf("invalid component id")
		return http.StatusBadRequest, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, fmt.Errorf("invalid component id")
	}

	_, prefix, err := query.GetComponentNameAndPrefix(ComponentID)
	if err != nil {
		log.Printf("error while fetching assigned units: %v", err)
		return http.StatusInternalServerError, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, fmt.Errorf("database error")
	}

	exists, err := query.CheckIfComponentBelongsToWarehouse(ComponentID, claims.UserID)
	if err != nil {
		log.Printf("error while fetching assigned units: %v", err)
		return http.StatusInternalServerError, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, fmt.Errorf("database error")
	} else if !exists {
		log.Printf("invalid component id")
		return http.StatusBadRequest, []models.AllOutOfWarentyWarehouseModel{}, -1, -1, -1, fmt.Errorf("component not found")
	}

	status, units, total, err := query.GetAllOutOfWarehouseUnitsInWarehouse(claims.UserID, Sort.Limit, Sort.Offset, prefix)
	if err != nil {
		log.Printf("error while fetching assigned units: %v", err)
		return http.StatusInternalServerError, []models.AllOutOfWarentyWarehouseModel{}, total, -1, -1, fmt.Errorf("database error")
	}

	return status, units, total, Sort.Limit, Sort.Page, nil
}
