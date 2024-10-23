package activities

import (
	"context"
	"errors"
	"idempotency/worker/internal/entities"
	"idempotency/worker/internal/repositories"

	"go.temporal.io/sdk/activity"
)

type CreateTransactionActivity interface {
	CreateTransaction(ctx context.Context, key, customerId string, amount float64) error
}

type createTransactionActivity struct {
	transactionRepo repositories.TransactionRepository
}

func NewCreateTransactionActivity(transactionRepo repositories.TransactionRepository) *createTransactionActivity {
	return &createTransactionActivity{
		transactionRepo: transactionRepo,
	}
}

func (a *createTransactionActivity) CreateTransaction(ctx context.Context, key, customerId string, amount float64) (string, error) {
	logger := activity.GetLogger(ctx)

	return "", errors.New("intentionally")

	logger.Info("create transaction started")

	var err error
	tx := a.transactionRepo.BeginTx(ctx)

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	transaction := entities.Transaction{
		UUID:       key,
		CustomerId: customerId,
		Amount:     amount,
	}

	if err = a.transactionRepo.CreateTransaction(ctx, transaction, tx); err != nil {
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	logger.Info("create transaction finished")

	return transaction.UUID, nil
}
