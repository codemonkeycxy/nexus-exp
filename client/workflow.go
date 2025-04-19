package client

import (
	"fmt"

	"go.temporal.io/sdk/workflow"

	greeting "nexus-exp/gen/proto/v1"
	"nexus-exp/gen/proto/v1/greetingnexus"
)

const (
	TaskQueue    = "my-caller-workflow-task-queue"
	endpointName = "my-nexus-endpoint-name"
)

func GreetWorkflow(ctx workflow.Context, message string) (string, error) {
	c := workflow.NewNexusClient(endpointName, greetingnexus.GreetingServiceName)

	fut := c.ExecuteOperation(ctx, greetingnexus.GreetingGreetOperationName, &greeting.GreetInput{
		Name: message,
	}, workflow.NexusOperationOptions{})

	var res greeting.GreetOutput
	if err := fut.Get(ctx, &res); err != nil {
		return "", err
	}

	return res.Greeting, nil
}

func SlothGreetWorkflow(ctx workflow.Context, message string) (string, error) {
	c := workflow.NewNexusClient(endpointName, greetingnexus.GreetingServiceName)

	start := workflow.Now(ctx)
	fut := c.ExecuteOperation(ctx, greetingnexus.GreetingSlothGreetOperationName, &greeting.GreetInput{
		Name: message,
	}, workflow.NexusOperationOptions{
		Summary: "🌿 < Hello Sloth > 🦥💤 ^__^",
	})

	var res greeting.GreetOutput
	if err := fut.Get(ctx, &res); err != nil {
		return "", err
	}

	return fmt.Sprintf("After %s, the sloth responded: %s", workflow.Now(ctx).Sub(start), res.Greeting), nil
}
