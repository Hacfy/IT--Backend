package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
		log.Printf("Error checking user details:", err)
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
		log.Printf("Error checking user details:", err)
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
		log.Printf("Error checking user details:", err)
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
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusCreated, nil
}

func (wr *WarehouseRepo) AssignUnits(e echo.Context) (int, error) {

	status, claims, err := utils.VerifyUserToken(e, "warehouses", wr.db)
	if err != nil {
		return status, err
	}

	query := database.NewDBinstance(wr.db)

	ok, err := query.VerifyUser(claims.UserEmail, "warehouses", claims.UserID)
	if err != nil {
		log.Printf("Error checking user details:", err)
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
		return http.StatusInternalServerError, fmt.Errorf("database error")
	}

	return http.StatusCreated, nil

}
