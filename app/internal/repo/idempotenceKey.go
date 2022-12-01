package repo

import (
	"github.com/google/uuid"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
)

type IdempotenceKeyRepo interface {
	IdempotenceKey
}

type IdempotenceKey interface {
	Get(key uuid.UUID) (idempotenceKey *domain.IdempotenceKey, err error)
	Create(idempotenceKey *domain.IdempotenceKey) error
}
