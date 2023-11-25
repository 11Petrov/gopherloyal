package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/11Petrov/gopherloyal/internal/config"
	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	"github.com/11Petrov/gopherloyal/internal/storage/postgre"
)

func ProcessOrderUpdates(ctx context.Context, cfg *config.Config, store *postgre.Database) {
	log := logger.FromContext(ctx)

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ordersToProcess, err := store.RetrieveNewOrders(ctx)
			if err != nil {
				log.Errorf("error getting new orders: %s", err)
				continue
			}

			for _, order := range ordersToProcess {
				status, accrual, err := getOrderAccrualInfo(ctx, cfg.AccrualAddress, order.Number)
				if err != nil {
					log.Errorf("error querying accrual service: %s", err)
					continue
				}

				if status == models.StatusInvalid || status == models.StatusProcessed {
					err := store.UpdateOrderStatusAndAccrual(ctx, order.Number, status, accrual)
					if err != nil {
						log.Errorf("error updating order status: %s", err)
					}
				}
			}
		}
	}
}

func getOrderAccrualInfo(ctx context.Context, accrualAddress, orderNumber string) (string, float64, error) {
	log := logger.FromContext(ctx)
	client := &http.Client{}
	req, err := http.NewRequest("GET", accrualAddress+"/api/orders/"+orderNumber, nil)
	if err != nil {
		log.Error("failed to create request:", err)
		return "", 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to send request:", err)
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		retryAfterDuration, err := strconv.Atoi(retryAfter)
		if err != nil {
			log.Errorf("failed to parse Retry-After header: %v", err)
			return "", 0, err
		}
		retryDuration := time.Duration(retryAfterDuration) * time.Second
		time.Sleep(retryDuration)

		log.Errorf("rate limited, retry after %d seconds", retryAfterDuration)
		return "", 0, err
	} else if resp.StatusCode != http.StatusOK {
		log.Errorf("accrual service responded with status: %d", resp.StatusCode)
		return "", 0, err
	}

	var response models.AccrualResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Error("failed to decode response:", err)
		return "", 0, err
	}

	return response.Status, response.Accrual, nil
}
