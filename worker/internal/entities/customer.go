package entities

type Customer struct {
	UUID    string  `db:"customer_uuid"`
	Name    string  `db:"name"`
	Balance float64 `db:"current_balance"`
}
