package postgre

import (
	"context"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	"github.com/11Petrov/gopherloyal/internal/storage/errors"
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

func (d *Database) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	log := logger.FromContext(ctx)

	conn, err := d.db.Acquire(ctx)
	if err != nil {
		log.Errorf("error acquiring database connection: %s", err)
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Errorf("error beginning transaction: %s", err)
		return err
	}
	defer tx.Rollback(ctx)

	var currentBalance float64
	err = tx.QueryRow(ctx, "SELECT current_balance FROM Users WHERE user_id = $1", userID).Scan(&currentBalance)
	if err != nil {
		log.Errorf("error retrieving user balance: %s", err)
		return err
	}

	if currentBalance < sum {
		return errors.ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, "UPDATE Users SET current_balance = current_balance - $1, withdrawn = withdrawn + $1 WHERE user_id = $2", sum, userID)
	if err != nil {
		log.Errorf("error withdrawing funds: %s", err)
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO Withdrawals (user_id, order_number, sum, processed_at) VALUES ($1, $2, $3, NOW())",
		userID, orderNumber, sum)
	if err != nil {
		log.Errorf("error recording withdrawal: %s", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("error committing transaction: %s", err)
		return err
	}

	log.Info("transaction committed successfully")
	return nil
}

func (d *Database) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawals, error) {
	log := logger.FromContext(ctx)

	var withdrawals []models.Withdrawals

	rows, err := d.db.Query(ctx, "SELECT order_number, sum, processed_at FROM Withdrawals WHERE user_id = $1 ORDER BY processed_at ASC", userID)
	if err != nil {
		log.Errorf("error retrieving user withdrawals: %s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal models.Withdrawals
		err := rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			log.Errorf("error scanning withdrawal rows: %s", err)
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	log.Info("user withdrawals successfully retrieved")

	return withdrawals, nil
}
