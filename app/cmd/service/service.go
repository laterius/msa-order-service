package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
	"github.com/laterius/service_architecture_hw3/app/internal/transport/client/dbrepo"
	"github.com/laterius/service_architecture_hw3/app/internal/transport/server/api"
	transport "github.com/laterius/service_architecture_hw3/app/internal/transport/server/http"
	_ "github.com/laterius/service_architecture_hw3/app/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
	"log"
)

func main() {
	var cfg domain.Config
	err := configor.New(&configor.Config{Silent: true}).Load(&cfg, "config/config.yaml", "./config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbrepo.Dsn(cfg.Db),
	}), &gorm.Config{
		Logger: dblogger.Default.LogMode(dblogger.Info),
	})
	if err != nil {
		panic(err)
	}

	newService := service.NewService(db)
	server := gin.Default()

	//Handlers
	server.POST("/order", api.CreateOrderHandler(newService))
	server.GET("/probe/live", transport.RespondOk())
	server.GET("/probe/ready", transport.RespondOk())

	err = server.Run(fmt.Sprintf(":%s", cfg.Http.Port))
	if err != nil {
		log.Fatalf("server start failed: %s", err)
	}
}
