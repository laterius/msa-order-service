package dbrepo

import (
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/repo"
	"gorm.io/gorm"
)

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) repo.OrderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) Get() (orders []*domain.Order, err error) {
	err = r.db.Find(&orders).Error
	if err != nil {
	}
	return
}

func (r *orderRepo) Create(order *domain.Order) error {
	err := r.db.Create(order).Error
	if err != nil {
		return err
	}
	return nil
}
