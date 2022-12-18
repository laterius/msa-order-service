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

// ReserveGoods sends request inventory service to reserve goods
func ReserveGoods(orderId uuid.UUID, goodIds []uuid.UUID) ([]uuid.UUID, error) {
	endpoint := fmt.Sprintf("%s/reserveGoods", service.StorageHost)
	data := map[string]interface{}{
		"order_id":  orderId,
		"goods_ids": goodIds,
	}

	body, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("internal service while reserving goods ")
	}

	return goodIds, nil
}
