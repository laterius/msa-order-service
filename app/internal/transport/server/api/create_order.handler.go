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
	"os"
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

		host, _ := os.Hostname()
		log.Println(host)

		log.Println("ORDER: creation order")
		o := service.CreateOrder()
		err := service.Store(o)

		if err != nil {
			log.Println("ORDER: creation failed")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    gin.H{},
			})

			return
		}

		log.Println("ORDER: created success")

		newSaga := saga.Saga{}
		newSaga.SetName("ORDER: start saga")

		newSaga.AddStep(saga.Step{
			Name: "reserve body",
			Func: func() error {
				log.Println("STORAGE: start reservation")
				_, err := inventory.ReserveGoods(o.Id, goodIds)

				if err != nil {
					log.Println("STORAGE: reservation failed")
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": err.Error(),
						"data":    gin.H{},
					})
				}

				log.Println("STORAGE: end reservation")
				return nil
			},
			Compensation: func() error {
				log.Println("STORAGE: cancel reservation")

				err := inventory.CancelGoodsReservation(o.Id)
				if err != nil {
					log.Println("STORAGE: FAIL cancel reservation")
					return err
				}

				log.Println("STORAGE: end cancel reservation")
				return nil
			},
		})

		newSaga.AddStep(saga.Step{
			Name: "reserve courier",
			Func: func() error {
				log.Println("DELIVERY: start courier reservation")
				err := shipment.ReserveCourier(o.Id)

				if err != nil {
					log.Println("DELIVERY: FAIL courier reservation")
					return err
				}

				log.Println("DELIVERY: end courier reservation.")
				return nil
			},
			Compensation: func() error {
				log.Println("DELIVERY: start cancel courier reservation")

				err := shipment.CancelCourierReservation(o.Id)
				if err != nil {
					log.Println("DELIVERY: FAIL cancel courier reservation")
					return err
				}

				log.Println("DELIVERY: end cancel courier reservation")
				return nil
			},
		})

		newSaga.AddStep(saga.Step{
			Name: "make payment",
			Func: func() error {

				log.Println("PAYMENT: start payment")
				err := payments.MakePayment(o.Id, amount)

				if err != nil {
					log.Println("PAYMENT: FAIL payment")
					return err
				}

				log.Println("PAYMENT: end payment")
				return nil
			},
			Compensation: func() error {
				log.Println("PAYMENT: cancel payment")

				err := payments.CancelPayment(o.Id)
				if err != nil {
					log.Println("PAYMENT: FAIL cancel payment")
					return err
				}

				log.Println("PAYMENT: canceled")

				return nil
			},
		})

		coordinator := saga.NewCoordinator(newSaga)
		err = coordinator.Commit()

		if err != nil {
			log.Println("ORDER: start cancelled order")

			err := service.Delete(o)
			if err != nil {
				log.Println(err.Error())
				log.Println("ORDER: FAIL cancelled order")
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "ORDER FAILED",
				"data": gin.H{
					"order_id": o.Id,
				},
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
