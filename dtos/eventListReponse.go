package dtos

import "eventom-backend/models"

type EventListResponse struct {
	Events   []*models.Event
	Metadata *EventListMetadata
}
