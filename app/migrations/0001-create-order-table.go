package mixtures

import (
	"github.com/ezn-go/mixture"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/satori/go.uuid"
	"time"
)

func init() {
	type Order struct {
		Id        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
		UserID    int       `json:"userId"`
		Status    int       `json:"status"`
		Amount    int       `json:"amount"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`
	}

	mx := &gormigrate.Migration{
		ID:       "0001",
		Migrate:  mixture.CreateTableM(&Order{}),
		Rollback: mixture.DropTableR(&Order{}),
	}

	mixture.Add(mixture.ForAnyEnv, mx)
}
