package client

import (
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/hashicorp/go-multierror"

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

	var greetings []string
	var multiErr error
	wg := workflow.NewWaitGroup(ctx)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		workflow.Go(ctx, func(ctx workflow.Context) {
			defer wg.Done()
			start := workflow.Now(ctx)
			fut := c.ExecuteOperation(ctx, greetingnexus.GreetingSlothGreetOperationName, &greeting.GreetInput{
				Name: message,
			}, workflow.NexusOperationOptions{
				ScheduleToCloseTimeout: 15 * time.Minute, // If sloth doesn't respond in 15 minutes, let it sleep.
				Summary:                "🌿 < Hello Sloth > 🦥💤 ^__^",
			})

			var res greeting.GreetOutput
			if err := fut.Get(ctx, &res); err != nil {
				multiErr = multierror.Append(multiErr, err)
			}
			greetings = append(greetings, fmt.Sprintf("After %s, the sloth responded: %s", workflow.Now(ctx).Sub(start), res.Greeting))
		})
	}

	wg.Wait(ctx)
	return strings.Join(greetings, "\n"), multiErr
}
