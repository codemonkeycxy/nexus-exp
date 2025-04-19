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

func SlothGreetWorkflow(ctx workflow.Context, input *greeting.SlothGreetInput) (*greeting.SlothGreetOutput, error) {
	var response *greeting.SlothGreetOutput
	if err := workflow.ExecuteChildWorkflow(ctx, SlothSleepAndGreetWorkflow, SlothSleepAndGreetWorkflowInput{
		GreetInput: input,
		CountDown:  5,
	}).Get(ctx, &response); err != nil {
		return response, err
	}

	return response, nil
}

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
