package server

import (
	"time"

	"go.temporal.io/sdk/workflow"

	greeting "nexus-exp/gen/proto/v1"
)

func SlothGreetWorkflow(ctx workflow.Context, input *greeting.GreetInput) (*greeting.GreetOutput, error) {
	if err := workflow.Sleep(ctx, 5*time.Second); err != nil {
		return nil, err
	}

	return &greeting.GreetOutput{
		Greeting: "Hello, " + input.Name,
	}, nil
}
