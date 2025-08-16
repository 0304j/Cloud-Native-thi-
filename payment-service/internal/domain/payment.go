package domain

import "time"

type Payment struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
