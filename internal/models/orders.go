package models

import "time"

const (
	StatusNew        string = "NEW"
	StatusProcessing string = "PROCESSING"
	StatusInvalid    string = "INVALID"
	StatusProcessed  string = "PROCESSED"
)

type Orders struct {
	UserID     int       `json:"-"`
	Number     string    `json:"order_number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
