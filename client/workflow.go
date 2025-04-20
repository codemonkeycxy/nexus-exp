package client

import (
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/hashicorp/go-multierror"

	greeting "nexus-exp/gen/proto/v1"
	"nexus-exp/gen/proto/v1/greetingnexustemporal"
)

const (
	TaskQueue    = "my-caller-workflow-task-queue"
	endpointName = "my-nexus-endpoint-name"
)

func GreetWorkflow(ctx workflow.Context, message string) (string, error) {
	c := greetingnexustemporal.NewGreetingNexusClient(endpointName)
	res, err := c.Greet(ctx, &greeting.GreetInput{
		Name: message,
	}, workflow.NexusOperationOptions{})
	if err != nil {
		return "", err
	}

	return res.Greeting, nil
}

func SlothGreetWorkflow(ctx workflow.Context, message string, slothNames []string) (string, error) {
	c := greetingnexustemporal.NewGreetingNexusClient(endpointName)

	var greetings []string
	var multiErr error
	wg := workflow.NewWaitGroup(ctx)
	for _, slothName := range slothNames {
		wg.Add(1)
		workflow.Go(ctx, func(ctx workflow.Context) {
			defer wg.Done()
			start := workflow.Now(ctx)
			res, err := c.SlothGreet(ctx, &greeting.SlothGreetInput{
				Greeting:  message,
				SlothName: slothName,
			}, workflow.NexusOperationOptions{
				ScheduleToCloseTimeout: 15 * time.Minute, // If sloth doesn't respond in 15 minutes, let it sleep.
				Summary:                "🌿 < Hello Sloth > 🦥💤 ^__^",
			})
			if err != nil {
				multiErr = multierror.Append(multiErr, err)
			} else {
				greetings = append(greetings, fmt.Sprintf("After %s, the sloth %s responded: %s", workflow.Now(ctx).Sub(start), slothName, res.Greeting))
			}
		})
	}

	wg.Wait(ctx)
	return strings.Join(greetings, "\n"), multiErr
}
