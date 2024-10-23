package entities

import "time"

type Transaction struct {
	UUID       string    `db:"transaction_uuid"`
	CustomerId string    `db:"customer_uuid"`
	Amount     float64   `db:"amount"`
	CreatedAt  time.Time `db:"created_at"`
}
