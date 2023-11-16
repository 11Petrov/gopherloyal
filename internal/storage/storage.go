package storage

import (
	"context"

	"github.com/11Petrov/gopherloyal/internal/models"
)

type Store interface {
	UserRegister(ctx context.Context, user *models.UserAuth) error
	UserLogin(ctx context.Context, user *models.UserAuth) error
}
