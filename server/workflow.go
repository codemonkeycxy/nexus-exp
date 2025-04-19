package server

import (
	"time"

	"go.temporal.io/sdk/workflow"

	greeting "nexus-exp/gen/proto/v1"
)

type (
	SlothSleepAndGreetWorkflowInput struct {
		GreetInput *greeting.GreetInput
		CountDown  int
	}
)

func SlothGreetWorkflow(ctx workflow.Context, input *greeting.GreetInput) (*greeting.GreetOutput, error) {
	var response *greeting.GreetOutput
	if err := workflow.ExecuteChildWorkflow(ctx, SlothSleepAndGreetWorkflow, SlothSleepAndGreetWorkflowInput{
		GreetInput: input,
		CountDown:  5,
	}).Get(ctx, &response); err != nil {
		return response, err
	}

	return response, nil
}

func SlothSleepAndGreetWorkflow(ctx workflow.Context, input SlothSleepAndGreetWorkflowInput) (*greeting.GreetOutput, error) {
	if input.CountDown < 1 {
		return &greeting.GreetOutput{
			Greeting: "Hello, " + input.GreetInput.Name,
		}, nil
	}

	if err := workflow.Sleep(ctx, time.Second); err != nil {
		return nil, err
	}

	input.CountDown--
	return nil, workflow.NewContinueAsNewError(ctx, SlothSleepAndGreetWorkflow, input)
}
