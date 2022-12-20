package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	s "github.com/laterius/service_architecture_hw3/app/internal/service"
	"net/http"
)

// GetOrderHandler handles request to get order by ID
func GetOrderHandler(service s.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("orderId")

		orderId, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "order id can't parse",
				"data":    gin.H{},
			})
			return
		}

		order, err := service.Get(orderId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "order not found",
				"data":    gin.H{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Order found",
			"data": gin.H{
				"order_id": order.Id,
			},
		})
	}
}
