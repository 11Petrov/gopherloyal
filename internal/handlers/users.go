package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/11Petrov/gopherloyal/internal/auth"
	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	storageErrors "github.com/11Petrov/gopherloyal/internal/storage/errors"
)

type users interface {
	UserRegister(ctx context.Context, user *models.UserAuth) error
	UserLogin(ctx context.Context, user *models.UserAuth) error
}

type usersHandler struct {
	store users
}

func NewUsersHandler(store users) *usersHandler {
	return &usersHandler{
		store: store,
	}
}

func (u *usersHandler) UserRegister(rw http.ResponseWriter, r *http.Request) {
	log := logger.LoggerFromContext(r.Context())

	var user models.UserAuth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Errorf("error decoding request body: %s", err)
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := u.store.UserRegister(r.Context(), &user)
	if err != nil {
		switch err {
		case storageErrors.ErrLoginTaken:
			http.Error(rw, "Login already taken", http.StatusConflict)
		default:
			log.Errorf("error registering user: %s", err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	token, err := auth.GenerateToken(r.Context(), user.Login)
	if err != nil {
		log.Errorf("error generating token: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Authorization", "Bearer "+token)
	rw.WriteHeader(http.StatusOK)
}

func (u *usersHandler) UserLogin(rw http.ResponseWriter, r *http.Request) {
	log := logger.LoggerFromContext(r.Context())

	var user models.UserAuth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Errorf("error decoding request body: %s", err)
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := u.store.UserLogin(r.Context(), &user)
	if err != nil {
		switch err {
		case storageErrors.ErrUserNotFound, storageErrors.ErrInvalidPassword:
			http.Error(rw, "Invalid login credentials", http.StatusUnauthorized)
		default:
			log.Errorf("error logging in user: %s", err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	token, err := auth.GenerateToken(r.Context(), user.Login)
	if err != nil {
		log.Errorf("error generating token: %s", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Authorization", "Bearer "+token)
	rw.WriteHeader(http.StatusOK)
}
