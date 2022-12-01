package dbrepo

import (
	"github.com/google/uuid"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/repo"
	"gorm.io/gorm"
)

type idempotenceKeyRepo struct {
	db *gorm.DB
}

func NewIdempotenceKeyRepo(db *gorm.DB) repo.IdempotenceKeyRepo {
	return &idempotenceKeyRepo{db: db}
}

func (r *idempotenceKeyRepo) Get(key uuid.UUID) (idempotenceKey *domain.IdempotenceKey, err error) {
	err = r.db.Model(idempotenceKey).Where("id = ?", key).First(&idempotenceKey).Error
	return
}

func (r *idempotenceKeyRepo) Create(idempotenceKey *domain.IdempotenceKey) error {
	err := r.db.Create(idempotenceKey).Error
	if err != nil {
		return err
	}
	return nil
}
