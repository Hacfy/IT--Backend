package models

type ComponentModel struct {
	ComponentID   int    `json:"component_id"`
	ComponentName string `json:"component_name"`
	Prefix        string `json:"prefix"`
	WarehouseID   int    `json:"warehouse_id"`
	A_at          int    `json:"a_at"`
}

type CreateComponentModel struct {
	ComponentName string `json:"component_name" validate:"required"`
}
