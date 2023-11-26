package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/11Petrov/gopherloyal/internal/auth"
	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	storageErrors "github.com/11Petrov/gopherloyal/internal/storage/errors"
	"github.com/11Petrov/gopherloyal/internal/utils"
)

type balance interface {
	GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error)
	Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawals, error)
}

type balanceHandler struct {
	store balance
}

func NewBalanceHandler(store balance) *balanceHandler {
	return &balanceHandler{
		store: store,
	}
}

func (u *balanceHandler) GetUserBalance(rw http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID, err := auth.GetUserID(r.Context(), r)
	if err != nil {
		log.Errorf("error getting user ID from token: %s", err)
		http.Error(rw, "User not authenticated", http.StatusUnauthorized)
		return
	}

	userBalance, err := u.store.GetUserBalance(r.Context(), userID)
	if err != nil {
		log.Errorf("error retrieving user balance: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(rw).Encode(userBalance); err != nil {
		log.Errorf("error encoding balance to JSON: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (u *balanceHandler) Withdrawals(rw http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID, err := auth.GetUserID(r.Context(), r)
	if err != nil {
		log.Errorf("error getting user ID from token: %s", err)
		http.Error(rw, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req models.Withdrawals
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("error decoding withdraw request: %s", err)
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	if !utils.IsLuhnValid(req.Order) {
		log.Warn("invalid order number format")
		http.Error(rw, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	err = u.store.Withdraw(r.Context(), userID, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, storageErrors.ErrInsufficientFunds):
			http.Error(rw, "Insufficient funds", http.StatusPaymentRequired)
			return
		default:
			log.Errorf("error withdrawing funds: %s", err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
}

func (u *balanceHandler) GetWithdrawals(rw http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userID, err := auth.GetUserID(r.Context(), r)
	if err != nil {
		log.Errorf("error getting user ID from token: %s", err)
		http.Error(rw, "User not authenticated", http.StatusUnauthorized)
		return
	}

	withdrawals, err := u.store.GetWithdrawals(r.Context(), userID)
	if err != nil {
		log.Errorf("error retrieving user withdrawals: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(rw).Encode(withdrawals); err != nil {
		log.Errorf("error encoding withdrawals to JSON: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}
