package models

import "github.com/labstack/echo/v4"

type DepartmentModel struct {
	DepartmentID int    `json:"department_id"`
	BranchID     int    `json:"branch_id"`
	BranchName   string `json:"branch_name"`
}

type CreateDepartmentModel struct {
	DepartmentName      string `json:"department_name" validate:"required"`
	DepartmentHeadName  string `json:"department_head_name" validate:"required"`
	DepartmentHeadEmail string `json:"department_head_email" validate:"required,email"`
}

type UpdateDepartmentHeadModel struct {
	DepartmentID           int    `json:"department_id" validate:"required"`
	DepartmentHeadEmail    string `json:"department_head_email" validate:"required,email"`
	NewDepartmentHeadName  string `json:"new_department_head_name" validate:"required"`
	NewDepartmentHeadEmail string `json:"new_department_head_email" validate:"required,email"`
}

type DepartmentInterface interface {
	CreateWorkspace(echo.Context) (int, int, error)
	DeleteWorkspace(echo.Context) (int, error)
	RaiseIssue(echo.Context) (int, int, error)
	RequestNewUnits(echo.Context) (int, map[int]int, error)
	GetAllDepartmentRequests(echo.Context) (int, []AllRequestsModel, error)
	GetDepartmentRequestDetails(echo.Context) (int, RequestDetailsModel, error)
	DeleteRequest(echo.Context) (int, error)
	DeleteIssue(echo.Context) (int, error)
}

type GetDepartmentIssuesModel struct {
	DepartmentID int `json:"department_id" validate:"required"`
}

type GetDepartmentWorkspacesModel struct {
	DepartmentID int `json:"department_id" validate:"required"`
}

type DeleteDepartmentModel struct {
	DepartmentID       int    `json:"department_id" validate:"required"`
	BranchHeadPassword string `json:"branch_head_password" validate:"required"`
}

type AllRequestsModel struct {
	RequestID     int    `json:"request_id"`
	WorkspaceID   int    `json:"workspace_id"`
	WarehouseID   int    `json:"warehouse_id"`
	DepartmentID  int    `json:"department_id"`
	ComponentID   int    `json:"component_id"`
	NumberOfUnits int    `json:"number_of_units"`
	Prefix        string `json:"prefix"`
	CreatedAt     string `json:"created_at"`
	Status        string `json:"status"`
}

type GetAllRequestsModel struct {
	DepartmentID int `json:"department_id" validate:"required"`
}

type RequestNewUnitModel struct {
	ComponentIDNoOfUnits map[int]int `json:"component_id_no_of_units" validate:"required"`
	WarehouseID          int         `json:"warehouse_id" validate:"required"`
	WorkspaceID          int         `json:"workspace_id" validate:"required"`
	DepartmentID         int         `json:"department_id" validate:"required"`
}

type RequestDetailsModel struct {
	RequestID     int    `json:"request_id"`
	WorkspaceID   int    `json:"workspace_id"`
	WarehouseID   int    `json:"warehouse_id"`
	DepartmentID  int    `json:"department_id"`
	ComponentID   int    `json:"component_id"`
	NumberOfUnits int    `json:"number_of_units"`
	Prefix        string `json:"prefix"`
	CreatedBy     int    `json:"created_by"`
	CreatedAt     string `json:"created_at"`
	Status        string `json:"status"`
}

type GetRequestDetailsModel struct {
	RequestID    int `json:"request_id" validate:"required"`
	DepartmentID int `json:"department_id" validate:"required"`
}

type DeleteRequestModel struct {
	RequestID int `json:"request_id" validate:"required"`
}
