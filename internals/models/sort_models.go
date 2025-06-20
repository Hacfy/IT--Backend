package models

type SortModel struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page" validate:"required"`
	Order  string `json:"order"`
	SortBy string `json:"sortBy"`
	Search string `json:"search"`
	Offset int    `json:"offset"`
}
