package repository

import (
	"database/sql"

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
	return -1, nil
}
