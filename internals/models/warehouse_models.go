package models

type WarehouseModel struct {
	WarehouseID        int    `json:"warehouse_id"`
	WarehouseUserName  string `json:"warehouse_user_name"`
	WarehouseUserEmail string `json:"warehouse_user_email"`
	BranchID           int    `json:"branch_id"`
}

type CreateWarehouseModel struct {
	WarehouseUserName  string `json:"warehouse_user_name" validate:"required"`
	WarehouseUserEmail string `json:"warehouse_user_email" `
}

type UpdateWarehouseHeadModel struct {
	WarehouseID           int    `json:"warehouse_id" validate:"required"`
	WarehouseHeadEmail    string `json:"warehouse_head_email" validate:"required,email"`
	NewWarehouseHeadName  string `json:"new_warehouse_head_name" validate:"required"`
	NewWarehouseHeadEmail string `json:"new_warehouse_head_email" validate:"required,email"`
}
