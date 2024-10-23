package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WithdrawArgs struct {
	CustomerId string  `json:"customerId"`
	Amount     float64 `json:"amount"`
}

func WithdrawCustomerWorkflow(ctx workflow.Context, args WithdrawArgs) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Hour,
	})

	ctx = workflow.WithRetryPolicy(ctx, temporal.RetryPolicy{
		MaximumAttempts: 4,
		MaximumInterval: 10 * time.Second,
	})

	key := workflow.GetInfo(ctx).WorkflowExecution.RunID
	err := workflow.ExecuteActivity(ctx, "WithdrawCustomer", key, args.CustomerId, args.Amount).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "CreateTransaction", key, args.CustomerId, args.Amount).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
