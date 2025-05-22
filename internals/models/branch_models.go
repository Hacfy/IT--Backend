package models

type BranchModel struct {
	BranchID       int    `json:"branch_id"`
	OrgID          int    `json:"org_id"`
	SuperAdminID   int    `json:"super_admin_id"`
	BranchName     string `json:"branch_name"`
	BranchLocation string `json:"branch_location"`
}

type BranchHeadModel struct {
	BranchID           int    `json:"branch_id"`
	BranchHeadID       int    `json:"branch_head_id"`
	BranchHeadName     string `json:"branch_head_name"`
	BranchHeadEmail    string `json:"branch_head_email"`
	BranchHeadPassword string `json:"branch_head_password"`
}

type CreateBranchModel struct {
	BranchName      string `json:"branch_name" validate:"required"`
	BranchLocation  string `json:"branch_location" validate:"required"`
	BranchHeadName  string `json:"branch_head_name" validate:"required"`
	BranchHeadEmail string `json:"branch_head_email" validate:"required,email"`
}
