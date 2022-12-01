package domain

import "github.com/google/uuid"

type IdempotenceKey struct {
	ID     uuid.UUID   `gorm:"not null;unique_index" json:"id"`
	Status OrderStatus `json:"status"`
}
