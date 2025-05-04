package server

import (
	"time"

	"go.temporal.io/sdk/workflow"

	greeting "nexus-exp/gen/proto/v1"
)

type (
	SlothSleepAndGreetWorkflowInput struct {
		GreetInput *greeting.SlothGreetInput
		CountDown  int
	}
)

func SlothSleepAndGreetWorkflow(ctx workflow.Context, input SlothSleepAndGreetWorkflowInput) (*greeting.SlothGreetOutput, error) {
	if input.CountDown < 1 {
		return &greeting.SlothGreetOutput{
			Greeting: "Hello, " + input.GreetInput.Greeting,
		}, nil
	}

	if err := workflow.Sleep(ctx, time.Second); err != nil {
		return nil, err
	}

	input.CountDown--
	return nil, workflow.NewContinueAsNewError(ctx, SlothSleepAndGreetWorkflow, input)
}
