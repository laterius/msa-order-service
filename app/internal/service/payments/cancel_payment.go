package payments

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
	"io"
	"log"
	"net/http"
)

// CancelPayment sends request to payment service to cancel payment
func CancelPayment(orderId uuid.UUID) error {
	endpoint := fmt.Sprintf("%s/cancelPayment", service.Host)
	data := map[string]interface{}{
		"order_id": orderId,
	}

	body, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return errors.New("internal service error")
	}

	return nil
}
