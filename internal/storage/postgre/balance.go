package postgre

import (
	"context"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
)

func (d *Database) GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error) {
	log := logger.FromContext(ctx)

	var userBalance models.UserBalance

	err := d.db.QueryRow(ctx, "SELECT current_balance, withdrawn FROM Users WHERE user_id = $1", userID).Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil {
		log.Errorf("error retrieving user balance: %s", err)
		return nil, err
	}

	log.Info("user balance successfully retrieved")

	return &userBalance, nil
}
