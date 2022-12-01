package domain

import (
	"github.com/google/uuid"
)

type OrderId uuid.UUID

type OrderStatus int

const (
	OrderStatusCreated OrderStatus = iota
)

type Order struct {
	ID     uuid.UUID   `gorm:"not null;unique_index" json:"id"`
	UserID int         `gorm:"not null;unique_index" json:"userId"`
	Status OrderStatus `json:"status"`
	Amount int         `json:"amount"`
}
