package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	s "github.com/laterius/service_architecture_hw3/app/internal/service"
	"io"
	"log"
	"net/http"
)

type Body struct {
	UserId uuid.UUID `json:"userId"`
	Amount int       `json:"amount"`
}

// CreateOrderHandlerV2 CreateOrderHandler handles request to create order
func CreateOrderHandlerV2() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := Body{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Problem with parse body request. Error = %s", err.Error()),
				"data":    gin.H{},
			})

			return
		}

		orderId := uuid.New()
		status := changeBalance(orderId, req, c)
		status = sendNotification(orderId, status, c)

		if status {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Order created successfully",
				"data": gin.H{
					"order_id": orderId,
				},
			})
		}
	}
}

func sendNotification(orderId uuid.UUID, status bool, c *gin.Context) bool {
	endpoint := fmt.Sprintf("%s/send", s.NotificationHost)
	data := map[string]interface{}{
		"order_id": orderId,
		"status":   status,
	}

	body, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Problem with send status from order. Error = %s", err.Error()),
			"data":    gin.H{},
		})
		return false
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Problem with send status from order. Status = %v", response.Status),
			"data":    gin.H{},
		})
		return false
	}
	return true
}

func changeBalance(orderId uuid.UUID, req Body, c *gin.Context) bool {
	endpoint := fmt.Sprintf("%s/balance/%s", s.BillingHost, req.UserId)
	data := map[string]interface{}{
		"amount": req.Amount,
	}

	body, _ := json.Marshal(data)

	request, _ := http.NewRequest("PUT", endpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Problem with change balance from order. Error = %s", err.Error()),
			"data":    gin.H{},
		})
		return false
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Problem with change balance from order. Status = %v", response.Status),
			"data": gin.H{
				"order_id": orderId,
			},
		})
		return false
	}

	return true
}
