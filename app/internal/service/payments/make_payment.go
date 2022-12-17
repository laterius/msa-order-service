package payments

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
)

// MakePayment sends request to payment service to make payment
func MakePayment(orderId uuid.UUID, amount int) error {
	endpoint := fmt.Sprintf("%s/makePayment", os.Getenv("PAYMENTS_HOST"))
	data := map[string]interface{}{
		"order_id": orderId,
		"amount":   amount,
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
