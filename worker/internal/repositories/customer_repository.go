package repositories

import (
	"context"
	"idempotency/worker/internal/entities"

	"github.com/jmoiron/sqlx"
)

type CustomerRepository interface {
	GetCustomerById(ctx context.Context, id string, tx Tx) (*entities.Customer, error)
	UpdateBalance(ctx context.Context, customerId string, newBalance float64, tx Tx) error
	BeginTx(ctx context.Context) Tx
}

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *customerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) GetCustomerById(ctx context.Context, id string, tx Tx) (*entities.Customer, error) {
	const sql = `select * from public.customers c where c.customer_uuid = $1;`
	var customer entities.Customer

	if err := tx.GetContext(ctx, &customer, sql, id); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) UpdateBalance(ctx context.Context, customerId string, newBalance float64, tx Tx) error {
	const sql = `update public.customers set current_balance = $1 where customer_uuid = $2;`
	_, err := tx.ExecContext(ctx, sql, newBalance, customerId)

	return err
}

func (r *customerRepository) BeginTx(ctx context.Context) Tx {
	return r.db.MustBeginTx(ctx, nil)
}
