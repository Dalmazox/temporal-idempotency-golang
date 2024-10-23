package activities

import (
	"context"
	"database/sql"
	"idempotency/worker/internal/repositories"

	"go.temporal.io/sdk/activity"
)

type WithdrawCustomerActivity interface {
	WithdrawCustomer(ctx context.Context, customerId string, amount float64) error
}

type withdrawCustomerActivity struct {
	customersRepo    repositories.CustomerRepository
	transactionsRepo repositories.TransactionRepository
}

func NewWithdrawCustomerActivity(customerRepo repositories.CustomerRepository, transactionsRepo repositories.TransactionRepository) *withdrawCustomerActivity {
	return &withdrawCustomerActivity{
		customersRepo:    customerRepo,
		transactionsRepo: transactionsRepo,
	}
}

func (a *withdrawCustomerActivity) WithdrawCustomer(ctx context.Context, key, customerId string, amount float64) error {
	logger := activity.GetLogger(ctx)

	logger.Info("withdraw customer start")

	var err error
	tx := a.customersRepo.BeginTx(ctx)

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	transaction, err := a.transactionsRepo.GetTransactionById(ctx, key, tx)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if transaction != nil {
		logger.Info("transaction already finished with uuid: %s", transaction.UUID)
		return nil
	}

	customer, err := a.customersRepo.GetCustomerById(ctx, customerId, tx)
	if err != nil {
		return err
	}

	newBalance := customer.Balance - amount
	if err = a.customersRepo.UpdateBalance(ctx, customer.UUID, newBalance, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return nil
	}

	logger.Info("withdraw customer finished")

	return nil
}
