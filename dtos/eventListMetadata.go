package dtos

type EventListMetadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}
