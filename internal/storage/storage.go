package storage

import (
	"context"

	"github.com/11Petrov/gopherloyal/internal/models"
)

type Store interface {
	UserRegister(ctx context.Context, user *models.Users) (int, error)
	UserLogin(ctx context.Context, user *models.Users) (*models.Users, error)
	UploadOrder(ctx context.Context, order *models.Orders) error
	GetUserOrders(ctx context.Context, userID int) ([]models.Orders, error)
	GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error)
}
