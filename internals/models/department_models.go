package models

type DepartmentModel struct {
	DepartmentID int    `json:"department_id"`
	BranchID     int    `json:"branch_id"`
	BranchName   string `json:"branch_name"`
}

type CreateDepartmentModel struct {
	DepartmentName      string `json:"department_id" validate:"required"`
	DepartmentHeadName  string `json:"department_head_name" validate:"required"`
	DepartmentHeadEmail string `json:"department_head_email" validate:"required,email"`
}

type UpdateDepartmentHeadModel struct {
	DepartmentID           int    `json:"department_id" validate:"required"`
	DepartmentHeadEmail    string `json:"department_head_email" validate:"required,email"`
	NewDepartmentHeadName  string `json:"new_department_head_name" validate:"required"`
	NewDepartmentHeadEmail string `json:"new_department_head_email" validate:"required,email"`
}
