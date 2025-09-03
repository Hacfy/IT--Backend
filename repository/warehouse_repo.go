package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/Hacfy/IT_INVENTORY/internals/models"
	"github.com/Hacfy/IT_INVENTORY/pkg/database"
	"github.com/Hacfy/IT_INVENTORY/pkg/utils"
	"github.com/labstack/echo/v4"
)

type WarehouseRepo struct {
	db *sql.DB
}

func NewWarehouseRepo(db *sql.DB) *WarehouseRepo {
	return &WarehouseRepo{
		db: db,
	}
}

func (wr *WarehouseRepo) generateUniquePrefix(name string) (string, error) {

	query := database.NewDBinstance(wr.db)

	if len(strings.TrimSpace(name)) == 0 {
		return "", fmt.Errorf("product name cannot be empty")
	}
	name = strings.ToLower(name)

	letters := ""
	seen := make(map[rune]bool)

	for _, ch := range name {
		if unicode.IsLetter(ch) && !seen[ch] {
			letters += string(ch)
			seen[ch] = true
		}
		if len(letters) >= 3 {
			break
		}
	}

	filler := "zyxwvutsrqponmlkjihgfedcba"

	for len(letters) < 3 {
		letters += string(filler[0])
	}

	prefix := letters
	maxAttempts := 10000
	attempts := 0

	for {
		if !query.IfPrefixExists(prefix) {
			break
		}

		attempts++
		if attempts > maxAttempts {
			return "", fmt.Errorf("failed to generate unique prefix")
		}

		runes := []rune(prefix)
		changed := false

		for i := len(runes) - 1; i >= 0; i-- {
			index := strings.IndexRune(filler, runes[i])
			if index == -1 {
				return "", fmt.Errorf("invalid character in prefix: %c", runes[i])
			}
			if index < len(filler)-1 {
				runes[i] = rune(filler[index+1])
				for j := i + 1; j < len(runes); j++ {
					runes[j] = rune(filler[0])
				}
				changed = true
				break
			}
		}

		if !changed {
			prefix = string(filler[0]) + string(runes)
		} else {
			prefix = string(runes)
		}
	}

	return prefix, nil

}

func (wr *WarehouseRepo) CreateComponent(e echo.Context) (int, string, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, "", err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, "", fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, "", fmt.Errorf("invalid user details")
	}

	var new_component models.CreateComponentModel

	if err := e.Bind(&new_component); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, "", fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_component); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, "", fmt.Errorf("failed to validate request")
	}

	if ok := query.IfComponentExists(new_component.ComponentName, claims.UserID); !ok {
		log.Printf("component %v already exists", new_component.ComponentName)
		return http.StatusConflict, "", fmt.Errorf("component already exists")
	}

	Prefix, err := wr.generateUniquePrefix(new_component.ComponentName)
	if err != nil {
		log.Printf("error while generating prefix: %v", err)
		return http.StatusInternalServerError, "", fmt.Errorf("failed to generate prefix, please try again later")
	}

	component_id, err := query.CreateComponent(new_component.ComponentName, Prefix, claims.UserID)
	if err != nil {
		log.Printf("error while storing Component data in DB: %v", err)
		return http.StatusInternalServerError, "", fmt.Errorf("unable to create component at the moment, please try again later")
	}

	token, err := utils.GenerateComponentToken(component_id, new_component.ComponentName, Prefix)
	if err != nil {
		log.Printf("error while generating token: %v", err)
		return http.StatusInternalServerError, "", fmt.Errorf("failed to generate token, please try again later")
	}

	return http.StatusOK, token, nil
}

func (wr *WarehouseRepo) DeleteComponent(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var del_component models.DeleteComponentModel

	if err := e.Bind(&del_component); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(del_component); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	if status, err := query.DeleteComponent(del_component, claims.UserID); err != nil {
		log.Printf("error while deleting the component %v: %v", del_component.ComponentID, err)
		return status, fmt.Errorf("error while deleting the component")
	}

	return http.StatusNoContent, nil
}

func (wr *WarehouseRepo) AddComponentUnits(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var new_component_unit models.AddUnitModel

	if err := e.Bind(&new_component_unit); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_component_unit); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	prefix, exists, err := query.CheckIfComponentIDExists(new_component_unit.ComponentID, claims.UserID)
	if err != nil {
		log.Printf("error while checking if component exists: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if exists {
		log.Printf("component with id %v does not exist", new_component_unit.ComponentID)
		return http.StatusBadRequest, fmt.Errorf("component with id %v does not exist", new_component_unit.ComponentID)
	}

	status, err = query.CreateComponentUnit(new_component_unit.Warenty_Date, float32(new_component_unit.Cost), prefix, claims.UserID, new_component_unit.Number_of_units, new_component_unit.ComponentID)
	if err != nil {
		log.Printf("error while creating units of %v: %v", new_component_unit.ComponentID, err)
		return status, fmt.Errorf("database error")
	}

	return status, nil
}

func (wr *WarehouseRepo) AssignUnits(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var new_unit models.AssignUnitModel

	if err := e.Bind(&new_unit); err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if err := validate.Struct(new_unit); err != nil {
		log.Printf("failed to validate request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("failed to validate request")
	}

	prefix, exists, err := query.CheckIfComponentIDExists(new_unit.ComponentID, claims.UserID)
	if err != nil {
		log.Printf("error while checking if component exists: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	if !exists {
		log.Printf("component with id %v does not exist", new_unit.ComponentID)
		return http.StatusBadRequest, fmt.Errorf("component with id %v does not exist", new_unit.ComponentID)
	}

	status, err = query.AssignUnitWorkspace(new_unit.WorkspaceID, new_unit.UnitIDs, prefix)
	if err != nil {
		log.Printf("error while assigning units to workspace: %v", err)
		return status, fmt.Errorf("database error")
	}

	return status, nil

}

func (wr *WarehouseRepo) GetAllWarehouseIssues(e echo.Context) (int, []models.IssueModel, int, int, int, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, []models.IssueModel{}, 0, 0, 0, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, []models.IssueModel{}, 0, 0, 0, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.IssueModel{}, 0, 0, 0, fmt.Errorf("invalid user details")
	}

	var Sort models.SortModel

	Sort.Limit, err = strconv.Atoi(e.QueryParam("limit"))
	if err != nil {
		return http.StatusBadRequest, []models.IssueModel{}, 0, 0, 0, fmt.Errorf("invalid request format")
	}

	if Sort.Limit <= 0 || Sort.Limit > 100 {
		Sort.Limit = 10
	}

	Sort.Page, err = strconv.Atoi(e.QueryParam("page"))
	if err != nil {
		return http.StatusBadRequest, []models.IssueModel{}, 0, 0, 0, fmt.Errorf("invalid request format")
	}

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

	allowed := map[string]bool{"name": true, "created_at": true}
	if !allowed[Sort.SortBy] {
		Sort.SortBy = "created_at"
	}

	Sort.Search = e.QueryParam("search")

	status, issues, total, err := query.GetAllIssues(claims.UserID, Sort)
	if err != nil {
		log.Printf("Error while fetching issues: %v", err)
		return http.StatusInternalServerError, []models.IssueModel{}, 0, 0, 0, fmt.Errorf("database error")
	}

	return status, issues, total, Sort.Page, Sort.Limit, nil

}

func (wr *WarehouseRepo) GetAllWarehouseComponents(e echo.Context) (int, []models.AllWarehouseComponentsModel, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, []models.AllWarehouseComponentsModel{}, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, []models.AllWarehouseComponentsModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.AllWarehouseComponentsModel{}, fmt.Errorf("invalid user details")
	}

	comps, err := query.GetAllWarehouseComponents(claims.UserID)
	if err != nil {
		log.Printf("error while fetching components: %v", err)
		return http.StatusInternalServerError, []models.AllWarehouseComponentsModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, comps, nil

}

func (wr *WarehouseRepo) GetAllWarehouseComponentUnits(e echo.Context) (int, []models.AllComponentUnitsModel, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, []models.AllComponentUnitsModel{}, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, []models.AllComponentUnitsModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.AllComponentUnitsModel{}, fmt.Errorf("invalid user details")
	}

	var getAllComponentUnitsModel models.GetAllComponentUnitsModel

	err = e.Bind(&getAllComponentUnitsModel)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, []models.AllComponentUnitsModel{}, fmt.Errorf("invalid request format")
	}

	if getAllComponentUnitsModel.ComponentID <= 0 {
		log.Printf("invalid component id")
		return http.StatusBadRequest, []models.AllComponentUnitsModel{}, fmt.Errorf("invalid component id")
	}

	units, err := query.GetAllWarehouseComponentUnits(getAllComponentUnitsModel.ComponentID)
	if err != nil {
		log.Printf("error while fetching components: %v", err)
		return http.StatusInternalServerError, []models.AllComponentUnitsModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, units, nil

}

func (wr *WarehouseRepo) GetIssueDetails(e echo.Context) (int, models.IssueDetailsModel, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, models.IssueDetailsModel{}, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, models.IssueDetailsModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, models.IssueDetailsModel{}, fmt.Errorf("invalid user details")
	}

	var getIssueDetailsModel models.GetIssueDetailsModel

	err = e.Bind(&getIssueDetailsModel)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, models.IssueDetailsModel{}, fmt.Errorf("invalid request format")
	}

	if getIssueDetailsModel.IssueID <= 0 {
		log.Printf("invalid issue id")
		return http.StatusBadRequest, models.IssueDetailsModel{}, fmt.Errorf("invalid issue id")
	}

	issue, err := query.GetIssueDetails(getIssueDetailsModel.IssueID)
	if err != nil {
		log.Printf("error while fetching issue details: %v", err)
		return http.StatusInternalServerError, models.IssueDetailsModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, issue, nil

}

func (wr *WarehouseRepo) GetUnitAssignmentHistory(e echo.Context) (int, models.UnitAssignmentHistoryModel, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, models.UnitAssignmentHistoryModel{}, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, models.UnitAssignmentHistoryModel{}, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, models.UnitAssignmentHistoryModel{}, fmt.Errorf("invalid user details")
	}

	var getUnitAssignmentHistoryModel models.GetUnitAssignmentHistoryModel

	err = e.Bind(&getUnitAssignmentHistoryModel)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, models.UnitAssignmentHistoryModel{}, fmt.Errorf("invalid request format")
	}

	if getUnitAssignmentHistoryModel.UnitID <= 0 {
		log.Printf("invalid unit id")
		return http.StatusBadRequest, models.UnitAssignmentHistoryModel{}, fmt.Errorf("invalid unit id")
	}

	history, err := query.GetUnitAssignmentHistory(getUnitAssignmentHistoryModel.UnitID)
	if err != nil {
		log.Printf("error while fetching unit assignment history: %v", err)
		return http.StatusInternalServerError, models.UnitAssignmentHistoryModel{}, fmt.Errorf("database error")
	}

	return http.StatusOK, models.UnitAssignmentHistoryModel{
		UnitID:  getUnitAssignmentHistoryModel.UnitID,
		History: history,
		Total:   len(history),
	}, nil

}

func (wr *WarehouseRepo) UpdateIssueStatus(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var updateIssueStatusModel models.UpdateIssueStatusModel

	err = e.Bind(&updateIssueStatusModel)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if updateIssueStatusModel.IssueID <= 0 {
		log.Printf("invalid issue id")
		return http.StatusBadRequest, fmt.Errorf("invalid issue id")
	}

	if updateIssueStatusModel.Status == "accepted" || updateIssueStatusModel.Status == "resolved" || updateIssueStatusModel.Status == "raised" {
		log.Printf("invalid issue status")
		return http.StatusBadRequest, fmt.Errorf("invalid issue status")
	}

	status, err = query.UpdateIssueStatus(updateIssueStatusModel.IssueID, updateIssueStatusModel.Status)
	if err != nil {
		log.Printf("error while updating issue status: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return status, nil

}

func (wr *WarehouseRepo) UpdateComponentName(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var updateComponentNameModel models.UpdateComponentNameModel

	err = e.Bind(&updateComponentNameModel)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if updateComponentNameModel.ComponentID <= 0 {
		log.Printf("invalid component id")
		return http.StatusBadRequest, fmt.Errorf("invalid component id")
	}

	if updateComponentNameModel.ComponentName == "" {
		log.Printf("invalid component name")
		return http.StatusBadRequest, fmt.Errorf("invalid component name")
	}

	status, err = query.UpdateComponentName(updateComponentNameModel.ComponentID, updateComponentNameModel.ComponentName)
	if err != nil {
		log.Printf("error while updating component name: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return status, nil

}

func (wr *WarehouseRepo) GetAssignedUnits(e echo.Context) (int, []models.AssignedUnitsModel, int, int, int, error) {

	page, err := strconv.Atoi(e.QueryParam("page"))
	if err != nil {
		return http.StatusBadRequest, []models.AssignedUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid request format")
	}

	limit, err := strconv.Atoi(e.QueryParam("limit"))
	if err != nil {
		return http.StatusBadRequest, []models.AssignedUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid request format")
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, []models.AssignedUnitsModel{}, -1, -1, -1, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, []models.AssignedUnitsModel{}, -1, -1, -1, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, []models.AssignedUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid user details")
	}

	var getAssignedUnitsModel models.GetAssignedUnitsModel

	err = e.Bind(&getAssignedUnitsModel)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, []models.AssignedUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid request format")
	}

	if getAssignedUnitsModel.ComponentID <= 0 {
		log.Printf("invalid component id")
		return http.StatusBadRequest, []models.AssignedUnitsModel{}, -1, -1, -1, fmt.Errorf("invalid component id")
	}

	_, prefix, err := query.GetComponentNameAndPrefix(getAssignedUnitsModel.ComponentID)
	if err != nil {
		log.Printf("error while fetching assigned units: %v", err)
		return http.StatusInternalServerError, []models.AssignedUnitsModel{}, -1, -1, -1, fmt.Errorf("database error")
	}

	units, total, err := query.GetAssignedUnits(prefix, getAssignedUnitsModel.WorkspaceID, limit, offset)
	if err != nil {
		log.Printf("error while fetching assigned units: %v", err)
		return http.StatusInternalServerError, []models.AssignedUnitsModel{}, total, -1, -1, fmt.Errorf("database error")
	}

	return http.StatusOK, units, total, limit, page, nil

}

func (wr *WarehouseRepo) UpdateMaintenanceCost(e echo.Context) (int, error) {
	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("Invalid user details")
		return http.StatusUnauthorized, fmt.Errorf("invalid user details")
	}

	var updateMaintenanceCostModel models.UpdateMaintenanceCostModel

	err = e.Bind(&updateMaintenanceCostModel)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request format")
	}

	if updateMaintenanceCostModel.UnitID <= 0 {
		log.Printf("invalid unit id")
		return http.StatusBadRequest, fmt.Errorf("invalid unit id")
	}

	prefix, ok, err := query.CheckIfUnitIDExists(updateMaintenanceCostModel.UnitID, updateMaintenanceCostModel.ComponentID, claims.UserID)
	if err != nil {
		log.Printf("error while checking if component exists: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	} else if !ok {
		log.Printf("component with id %v does not exist", updateMaintenanceCostModel.UnitID)
		return http.StatusBadRequest, fmt.Errorf("component with id %v does not exist", updateMaintenanceCostModel.UnitID)
	}

	if updateMaintenanceCostModel.MaintenanceCost < 0 {
		log.Printf("invalid maintenance cost")
		return http.StatusBadRequest, fmt.Errorf("invalid maintenance cost")
	}

	status, err = query.UpdateMaintenanceCost(updateMaintenanceCostModel.UnitID, prefix, updateMaintenanceCostModel.MaintenanceCost)
	if err != nil {
		log.Printf("error while updating component name: %v", err)
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusOK, nil

}
