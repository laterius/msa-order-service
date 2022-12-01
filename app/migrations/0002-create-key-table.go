package mixtures

import (
	"github.com/ezn-go/mixture"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
)

func init() {
	type IdempotenceKey struct {
		Id     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
		Status int       `json:"status"`
	}

	mx := &gormigrate.Migration{
		ID:       "0002",
		Migrate:  mixture.CreateTableM(&IdempotenceKey{}),
		Rollback: mixture.DropTableR(&IdempotenceKey{}),
	}

	mixture.Add(mixture.ForAnyEnv, mx)
}
