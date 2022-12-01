package repo

import (
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
)

type OrderRepo interface {
	OrderReader
	OrderCreator
}

type OrderReader interface {
	Get() ([]*domain.Order, error)
}

type OrderCreator interface {
	Create(order *domain.Order) error
}
