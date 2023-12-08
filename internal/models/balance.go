package models

import "time"

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawals struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
