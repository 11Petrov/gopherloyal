package postgre

import (
	"context"
	"errors"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	"github.com/jackc/pgx/v5"
)

func (d *Database) RetrieveNewOrders(ctx context.Context) ([]models.Orders, error) {
	log := logger.FromContext(ctx)
	var orders []models.Orders

	query := `
        SELECT user_id, order_number, status, accrual, uploaded_at
        FROM Orders
        WHERE status = 'NEW'
    `

	rows, err := d.db.Query(ctx, query)
	if err != nil {
		log.Errorf("error retrieving new orders: %s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Orders
		if err := rows.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			log.Errorf("error scanning new order: %s", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Errorf("error during rows iteration for new orders: %s", err)
		return nil, err
	}

	return orders, nil
}

func (d *Database) UpdateOrderStatusAndAccrual(ctx context.Context, orderNumber string, status string, accrual float64) error {
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

	cmdTag, err := tx.Exec(ctx, `
		UPDATE Orders
		SET status = $1, accrual = $2
		WHERE order_number = $3
	`, status, accrual, orderNumber)
	if err != nil {
		log.Errorf("error executing order update: %s", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		log.Warn("no order found to update")
		return errors.New("no rows updated")
	}

	if status == "PROCESSED" {
		err = d.updateUserBalance(ctx, tx, orderNumber, accrual)
		if err != nil {
			log.Errorf("error updating user balance: %s", err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("error committing transaction: %s", err)
		return err
	}

	log.Info("order status and accrual successfully updated")
	return nil
}

func (d *Database) updateUserBalance(ctx context.Context, tx pgx.Tx, orderNumber string, accrual float64) error {
	var userID int
	err := tx.QueryRow(ctx, `
		SELECT user_id FROM Orders WHERE order_number = $1
	`, orderNumber).Scan(&userID)
	if err != nil {
		return err
	}

	cmdTag, err := tx.Exec(ctx, `
		UPDATE Users SET current_balance = current_balance + $2 WHERE user_id = $1
	`, userID, accrual)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}

	return nil
}
