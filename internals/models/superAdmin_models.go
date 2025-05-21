package models

type SuperAdminModel struct {
	ID       int    `json:"super_admin_id"`
	Org_ID   int    `json:"org_id"`
	Name     string `json:"super_admin_name"`
	Email    string `json:"super_admin_email"`
	Password string `json:"super_admin_password"`
}

type CreateSuperAdminModel struct {
	SuperAdminName  string `json:"super_admin_name" validate:"required"`
	SuperAdminEmail string `json:"super_admin_email" validate:"required,email"`
}

//login in auth model

type DeleteSuperAdminModel struct {
	SuperAdminEmail string `json:"super_admin_email" validate:"required,email"`
}
