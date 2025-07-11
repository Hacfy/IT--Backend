package models

import "time"

// type UnitModel struct {
// 	ID           string    `json:"id,omitempty"`
// 	Com_ID         *int      `json:"c_id" validate:"required"`
// 	U_ID         string    `json:"u_id,omitempty"`
// 	Assigned     bool      `json:"assigned"`
// 	W_ID         *int      `json:"w_id" validate:"required"`
// 	DW_ID        *int      `json:"dw_id,omitempty"`
// 	Warenty_Date time.Time `json:"warenty_date" validate:"required"`
// 	Status       string    `json:"status" validate:"required"`
// 	C_By         *int      `json:"c_by" validate:"required"`
// 	Cost         float64   `json:"cost" validate:"required"`
// 	M_Cost       float64   `json:"m_cost,omitempty"`
// }

type UnitModel struct {
	UnitID            int       `json:"unit_id"`
	ComponentID       int       `json:"component_id"`
	WarehouseID       int       `json:"warehouse_id"`
	WorkspaceID       *int      `json:"workspace_id"`
	Warenty_Date      time.Time `json:"warenty_date"`
	Status            string    `json:"status"`
	Cost              float64   `json:"cost"`
	Maintainance_Cost float64   `json:"maintainance_cost"`
	Created_By        int       `json:"created_by"`
}

type AddUnitModel struct {
	Number_of_units int       `json:"number_of_units" validate:"required"`
	ComponentID     int       `json:"component_id" validate:"required"`
	Warenty_Date    time.Time `json:"warenty_date" validate:"required"`
	Cost            float64   `json:"cost" validate:"required"`
}

type AssignUnitModel struct {
	WorkspaceID int    `json:"workspace_id" validate:"required"`
	ComponentID int    `json:"component_id" validate:"required"`
	UnitIDs     []int  `json:"unit_ids" validate:"required"`
	Prefix      string `json:"prefix" validate:"required"`
}
