package mixtures

import (
	"github.com/ezn-go/mixture"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func init() {

	mx := &gormigrate.Migration{
		ID:       "0001",
		Migrate:  mixture.CreateTableM(&service.Order{}),
		Rollback: mixture.DropTableR(&service.Order{}),
	}

	mixture.Add(mixture.ForAnyEnv, mx)
}
