package models

type Users struct {
	ID             int     `json:"-"`
	Login          string  `json:"login,omitempty"`
	Password       string  `json:"password,omitempty"`
	CurrentBalance float64 `json:"current_balance"`
	Withdrawn      float64 `json:"withdrawn"`
}
