package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/11Petrov/gopherloyal/internal/auth"
	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	storageErrors "github.com/11Petrov/gopherloyal/internal/storage/errors"
	"github.com/11Petrov/gopherloyal/internal/utils"
)

type orders interface {
	UploadOrder(ctx context.Context, order *models.Orders) error
	GetUserOrders(ctx context.Context, userID int) ([]models.Orders, error)
}

type ordersHandler struct {
	store orders
}

func NewOrdersHandler(store orders) *ordersHandler {
	return &ordersHandler{
		store: store,
	}
}

func (o *ordersHandler) UploadOrder(rw http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		log.Errorf("user not authenticated")
		http.Error(rw, "User not authenticated", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}
	orderNumber := string(body)

	if !utils.IsLuhnValid(orderNumber) {
		log.Warn("invalid order number format")
		http.Error(rw, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	order := &models.Orders{
		UserID:     userID,
		Number:     orderNumber,
		Status:     models.StatusNew,
		UploadedAt: time.Now(),
	}

	err = o.store.UploadOrder(r.Context(), order)
	if err != nil {
		switch {
		case errors.Is(err, storageErrors.ErrUploadedByThisUser):
			http.Error(rw, "Order has already been uploaded by this user", http.StatusOK)
			return
		case errors.Is(err, storageErrors.ErrUploadedByAnotherUser):
			http.Error(rw, "Order has already ben uploaded by another user", http.StatusConflict)
			return
		default:
			log.Errorf("error uploading order: %s", err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		rw.WriteHeader(http.StatusAccepted)
		return
	}
}

func (o *ordersHandler) GetUserOrders(rw http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		log.Errorf("user not authenticated")
		http.Error(rw, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Получаем список заказов пользователя из базы данных
	orders, err := o.store.GetUserOrders(r.Context(), userID)
	if err != nil {
		log.Errorf("error fetching user orders: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Возвращаем список заказов в формате JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(orders); err != nil {
		log.Errorf("error encoding user orders to JSON: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
