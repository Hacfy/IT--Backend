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

type DepartmentIssuesModel struct {
	IssueID     int    `json:"issue_id"`
	Issue       string `json:"issue"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UnitID      int    `json:"unit_id"`
	UnitPrefix  string `json:"unit_prefix"`
	WorkspaceID int    `json:"workspace_id"`
}

type GetIssueDetailsModel struct {
	IssueID int `json:"issue_id" validate:"required"`
}

type IssueDetailsModel struct {
	IssueID      int       `json:"issue_id"`
	DepartmentID int       `json:"department_id"`
	WarehouseID  int       `json:"warehouse_id"`
	WorkspaceID  int       `json:"workspace_id"`
	UnitID       int       `json:"unit_id"`
	UnitPrefix   string    `json:"unit_prefix"`
	Issue        string    `json:"issue"`
	Created_at   time.Time `json:"created_at"`
	Status       string    `json:"status"`
}

type UpdateIssueStatusModel struct {
	IssueID int    `json:"issue_id" validate:"required"`
	Status  string `json:"status" validate:"required"`
}

type DeleteIssueModel struct {
	IssueID int `json:"issue_id" validate:"required"`
}
