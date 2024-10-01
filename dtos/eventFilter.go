package dtos

type EventFilterDto struct {
	Page         int `validate:"required,gte=1"`
	PageSize     int `validate:"required,oneof=10 15 20 25"`
	Name         string
	Location     string
	FreeCapacity int    `validate:"omitempty,number,gte=0"`
	SortColumn   string `validate:"omitempty,oneof=id event_name event_description event_date max_capacity amount_registrations"`
	SortOrder    string `validate:"omitempty,oneof=DESC ASC"`
}
