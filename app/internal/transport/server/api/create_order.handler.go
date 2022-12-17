package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/laterius/service_architecture_hw3/app/internal/saga"
	s "github.com/laterius/service_architecture_hw3/app/internal/service"
	"github.com/laterius/service_architecture_hw3/app/internal/service/inventory"
	"github.com/laterius/service_architecture_hw3/app/internal/service/payments"
	"github.com/laterius/service_architecture_hw3/app/internal/service/shipment"
	"log"
	"net/http"
)

// CreateOrderHandler handles request to create order
func CreateOrderHandler(service s.Service) func(c *gin.Context) {
	type Body struct {
		Goods []s.Good
	}

	return func(c *gin.Context) {

		body := Body{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    gin.H{},
			})

			return
		}

		var goodIds []uuid.UUID
		amount := 0
		for _, good := range body.Goods {
			amount += good.Amount
			goodIds = append(goodIds, good.ID)
		}

		log.Println("create order")
		o := service.CreateOrder()
		err := service.Store(o)

		if err != nil {
			log.Println("order creation failed")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    gin.H{},
			})

			return
		}

		log.Println("order created")

		newSaga := saga.Saga{}
		newSaga.SetName("order creation")

		newSaga.AddStep(saga.Step{
			Name: "reserve body",
			Func: func() error {
				log.Println("inventory: start body reservation")
				_, err := inventory.ReserveGoods(o.Id, goodIds)

				if err != nil {
					log.Println("inventory reservation failed")
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": err.Error(),
						"data":    gin.H{},
					})
				}

				log.Println("inventory: end body reservation.")
				return nil
			},
			Compensation: func() error {
				log.Println("inventory: cancel body reservation")

				err := inventory.CancelGoodsReservation(o.Id)
				if err != nil {
					return err
				}

				return nil
			},
		})

		newSaga.AddStep(saga.Step{
			Name: "reserve courier",
			Func: func() error {
				log.Println("shipment: start courier reservation")
				err := shipment.ReserveCourier(o.Id)

				if err != nil {
					return err
				}

				log.Println("shipment: end courier reservation.")
				return nil
			},
			Compensation: func() error {
				log.Println("shipment: cancel courier reservation")

				err := shipment.CancelCourierReservation(o.Id)
				if err != nil {
					return err
				}

				return nil
			},
		})

		newSaga.AddStep(saga.Step{
			Name: "make payment",
			Func: func() error {

				log.Println("payments: start payment")
				err := payments.MakePayment(o.Id, amount)

				if err != nil {
					return err
				}

				log.Println("payments: end payment")
				return nil
			},
			Compensation: func() error {
				log.Println("payments: cancel payment")

				err := payments.CancelPayment(o.Id)
				if err != nil {
					return err
				}

				log.Println("payments: canceled")

				return nil
			},
		})

		coordinator := saga.NewCoordinator(newSaga)
		err = coordinator.Commit()

		if err != nil {
			log.Println("order cancelled")

			err := service.Delete(o)
			if err != nil {
				log.Println(err.Error())
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    gin.H{},
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data": gin.H{
				"order_id": o.Id,
			},
		})
	}
}
