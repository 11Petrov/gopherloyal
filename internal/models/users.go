package models

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}