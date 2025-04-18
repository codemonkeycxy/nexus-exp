package client

import (
	"go.temporal.io/sdk/workflow"

	greeting "nexus-exp/gen/proto/v1"
	"nexus-exp/gen/proto/v1/greetingnexus"
)

const (
	TaskQueue    = "my-caller-workflow-task-queue"
	endpointName = "my-nexus-endpoint-name"
)

func SlothGreetWorkflow(ctx workflow.Context, message string) (string, error) {
	c := workflow.NewNexusClient(endpointName, greetingnexus.GreetingServiceName)

	fut := c.ExecuteOperation(ctx, greetingnexus.GreetingSlothGreetOperationName, &greeting.GreetInput{
		Name: message,
	}, workflow.NexusOperationOptions{})

	var res greeting.GreetOutput
	if err := fut.Get(ctx, &res); err != nil {
		return "", err
	}

	return res.Greeting, nil
}
