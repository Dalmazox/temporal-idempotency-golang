package main

import (
	"idempotency/worker/config"
	"idempotency/worker/internal/activities"
	"idempotency/worker/internal/repositories"
	"idempotency/worker/internal/workflows"
	"log"

	_ "github.com/lib/pq"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("could not load the .env file")
	}

	db, err := getDatabase(cfg.Database.ConnectionString)
	if err != nil {
		log.Fatal("could not connect to database")
	}
	defer db.Close()

	c, err := getTemporalClient(cfg.Temporal.Address)
	if err != nil {
		log.Fatal("could not connect to temporal database")
	}
	defer c.Close()

	customerRepo := repositories.NewCustomerRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	withdrawCustomerActivity := activities.NewWithdrawCustomerActivity(customerRepo, transactionRepo)
	createTransactionActivity := activities.NewCreateTransactionActivity(transactionRepo)

	w := worker.New(c, "withdraw-customer", worker.Options{})

	w.RegisterWorkflow(workflows.WithdrawCustomerWorkflow)
	w.RegisterActivity(withdrawCustomerActivity.WithdrawCustomer)
	w.RegisterActivity(createTransactionActivity.CreateTransaction)

	if err = w.Run(worker.InterruptCh()); err != nil {
		log.Fatal("could not start the withdraw customer worker")
	}
}

func getDatabase(connectionString string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", connectionString)
}

func getTemporalClient(address string) (client.Client, error) {
	return client.Dial(client.Options{HostPort: address})
}
