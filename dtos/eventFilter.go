package dtos

type EventFilterDto struct {
	Name         string `validate:"omitempty,alphanum"`
	Location     string `validate:"omitempty,alpha"`
	FreeCapacity int    `validate:"omitempty,number,gte=0"`
	SortColumn   string `validate:"omitempty,oneof=id event_name event_description event_date max_capacity amount_registrations"`
	SortOrder    string `validate:"omitempty,oneof=DESC ASC"`
}
