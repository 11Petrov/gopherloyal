package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/11Petrov/gopherloyal/internal/auth"
	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
)

type balance interface {
	GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error)
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
	log := logger.LoggerFromContext(r.Context())

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
