package postgre

import (
	"context"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	storageErrors "github.com/11Petrov/gopherloyal/internal/storage/errors"
)

func (d *Database) UploadOrder(ctx context.Context, order *models.Orders) error {
	log := logger.FromContext(ctx)

	rows, err := d.db.Query(ctx, "SELECT user_id FROM Orders WHERE order_number = $1", order.Number)
	if err != nil {
		log.Errorf("error executing select query: %s", err)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var existingUserID int
		if err = rows.Scan(&existingUserID); err != nil {
			log.Errorf("error scanning user_id: %s", err)
			return err
		}
		if existingUserID == order.UserID {
			return storageErrors.ErrUploadedByThisUser
		}
		return storageErrors.ErrUploadedByAnotherUser
	}

	_, err = d.db.Exec(ctx, "INSERT INTO Orders (user_id, order_number, status, uploaded_at) VALUES ($1, $2, $3, $4)",
		order.UserID, order.Number, order.Status, order.UploadedAt)
	if err != nil {
		log.Errorf("error uploading order: %s", err)
		return err
	}
	return nil
}

func (d *Database) GetUserOrders(ctx context.Context, userID int) ([]models.Orders, error) {
	log := logger.FromContext(ctx)

	rows, err := d.db.Query(ctx, "SELECT order_number, uploaded_at, status, accrual FROM Orders WHERE user_id = $1 ORDER BY uploaded_at ASC", userID)
	if err != nil {
		log.Errorf("error fetching user orders: %s", err)
		return nil, err
	}
	defer rows.Close()

	var orders []models.Orders
	for rows.Next() {
		var order models.Orders
		if err := rows.Scan(&order.Number, &order.UploadedAt, &order.Status, &order.Accrual); err != nil {
			log.Errorf("error scanning order rows: %s", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Errorf("error scanning order rows: %s", err)
		return nil, err
	}

	return orders, nil
}
