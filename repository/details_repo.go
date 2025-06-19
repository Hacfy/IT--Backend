package repository

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/Hacfy/IT_INVENTORY/internals/models"
// 	"github.com/Hacfy/IT_INVENTORY/pkg/database"
// 	"github.com/labstack/echo/v4"
// )

// type DetaisRepo struct {
// 	db *sql.DB
// }

// func NewRepo(db *sql.DB) *DetaisRepo {
// 	return &DetaisRepo{db: db}
// }

// func (dr *DetaisRepo) GetAllDepartmentsRepo(e echo.Context, db *sql.DB) ([]models.GetAllDepartmentsModel, int, error) {
// 	limit, _ := strconv.Atoi(e.QueryParam("limit"))
// 	if limit <= 0 || limit > 100 {
// 		limit = 10
// 	}
// 	page, _ := strconv.Atoi(e.QueryParam("page"))
// 	if page <= 0 {
// 		page = 1
// 	}
// 	offset := (page - 1) * limit

// 	order := e.QueryParam("order")
// 	if order != "asc" && order != "desc" {
// 		order = "asc"
// 	}

// 	sortBy := e.QueryParam("sortBy")
// 	if sortBy == "" {
// 		sortBy = "created_at"
// 	}

// 	allowed := map[string]bool{"name": true, "created_at": true}
// 	if !allowed[sortBy] {
// 		sortBy = "created_at"
// 	}

// 	search := e.QueryParam("search")
// 	role := e.Get("userType").(string)

// 	// whereClause := ""
// 	// args := []interface{}{limit, offset}
// 	// argIndex := 3
// 	// if role == "admin" {
// 	// 	whereClause = "WHERE role = 'user'"
// 	// } else if role == "superAdmin" {
// 	// 	whereClause = "WHERE role IN ('admin', 'user')"
// 	// } else {
// 	// 	return []models.GetAllDepartmentsModel{}, http.StatusForbidden, fmt.Errorf("invalid creadentials")
// 	// }
// 	// if search != "" {
// 	// 	whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex)
// 	// 	args = append(args, "%"+search+"%")
// 	// }

// 	query := database.NewDBinstance(dr.db)

// }
