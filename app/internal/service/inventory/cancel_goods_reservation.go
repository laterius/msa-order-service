package inventory

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

// CancelGoodsReservation sends request to cancel reservation
func CancelGoodsReservation(orderId uuid.UUID) error {
	endpoint := fmt.Sprintf("%s/cancelGoodsReservation", service.StorageHost)
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
		return errors.New("internal service while canceling reservation goods")
	}

	return nil
}
