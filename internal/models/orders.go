package models

import "time"

type Orders struct {
	UserID      int       `json:"-"`
	OrderNumber string    `json:"order_number"`
	Status      string    `json:"status"`
	Accrual     float64   `json:"accrual,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
}