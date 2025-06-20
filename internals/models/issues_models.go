package models

import "time"

type IssueModel struct {
	IssueID      int       `json:"issue_id"`
	DepartmentID int       `json:"department_id" validate:"required"`
	WarehouseID  int       `json:"warehouse_id" validate:"required"`
	WorkspaceID  int       `json:"workspace_id" validate:"required"`
	UnitID       int       `json:"unit_id" validate:"required"`
	UnitPrefix   string    `json:"unit_prefix" validate:"required"`
	Issue        string    `json:"issue" validate:"required"`
	Created_at   time.Time `json:"created_at"`
	Status       string    `json:"status"`
}
