package models

import "github.com/golang-jwt/jwt/v5"

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

type ComponentTokenModel struct {
	ComponentID   int    `json:"component_id"`
	ComponentName string `json:"component_name"`
	Prefix        string `json:"prefix"`
	jwt.RegisteredClaims
}

type DeleteComponentModel struct {
	ComponentID   int    `json:"component_id" validate:"required"`
	ComponentName string `json:"component_name" validate:"required"`
	Prefix        string `json:"prefix" validate:"required"`
}
