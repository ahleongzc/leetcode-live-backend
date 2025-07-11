package model

type Pagination struct {
	Offset  uint `json:"offset"`
	Limit   uint `json:"limit"`
	Total   uint `json:"total"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}
