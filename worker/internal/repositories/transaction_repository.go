package repositories

import (
	"context"
	"idempotency/worker/internal/entities"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	GetTransactionById(ctx context.Context, id string, tx Tx) (*entities.Transaction, error)
	CreateTransaction(ctx context.Context, transaction entities.Transaction, tx Tx) error
	BeginTx(ctx context.Context) Tx
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetTransactionById(ctx context.Context, id string, tx Tx) (*entities.Transaction, error) {
	const sql = `select * from public.transactions t where t.transaction_uuid = $1;`
	var transaction entities.Transaction

	if err := tx.GetContext(ctx, &transaction, sql, id); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction entities.Transaction, tx Tx) error {
	const sql = `insert into public.transactions (transaction_uuid, customer_uuid, amount) values ($1, $2, $3)`

	_, err := tx.ExecContext(
		ctx,
		sql,
		transaction.UUID,
		transaction.CustomerId,
		transaction.Amount)

	return err
}

func (r *transactionRepository) BeginTx(ctx context.Context) Tx {
	return r.db.MustBeginTx(ctx, nil)
}
