package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/jinzhu/configor"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"

	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
	"github.com/laterius/service_architecture_hw3/app/internal/transport/client/dbrepo"
	transport "github.com/laterius/service_architecture_hw3/app/internal/transport/server/http"
	_ "github.com/laterius/service_architecture_hw3/app/migrations"
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

	orderRepo := dbrepo.NewOrderRepo(db)
	keyRepo := dbrepo.NewIdempotenceKeyRepo(db)
	orderService := service.NewOrderService(orderRepo, keyRepo)
	getOrderHandler := transport.NewGetOrder(orderService)
	postOrderHandler := transport.NewPostOrder(orderService)

	engine := html.New("./views", ".html")
	srv := fiber.New(fiber.Config{Views: engine})

	srv.Use(logger.New())
	srv.Use(favicon.New())
	srv.Use(recover.New())

	api := srv.Group("/api")
	api.Post("/order", postOrderHandler.Handle())
	api.Get("/orders", getOrderHandler.Handle())

	srv.Get("/probe/live", transport.RespondOk)
	srv.Get("/probe/ready", transport.RespondOk)

	srv.All("/*", transport.DefaultResponse)

	err = srv.Listen(fmt.Sprintf(":%s", cfg.Http.Port))
	if err != nil {
		panic(err)
	}
}
