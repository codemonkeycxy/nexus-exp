package main

import (
	"context"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporalnexus"
	"go.temporal.io/sdk/worker"

	greeting "nexus-exp/gen/proto/v1"
	"nexus-exp/gen/proto/v1/greetingnexus"
	"nexus-exp/options"
	"nexus-exp/server"

	"github.com/nexus-rpc/sdk-go/nexus"
)

const (
	taskQueue = "my-handler-task-queue"
)

type handler struct {
	greetingnexus.GreetingNexusHandler
}

func (h *handler) Greet(name string) nexus.Operation[*greeting.GreetInput, *greeting.GreetOutput] {
	return nexus.NewSyncOperation(name, func(ctx context.Context, input *greeting.GreetInput, options nexus.StartOperationOptions) (*greeting.GreetOutput, error) {
		return &greeting.GreetOutput{
			Greeting: "Hello, " + input.Name,
		}, nil
	})
}

func (h *handler) SlothGreet(name string) nexus.Operation[*greeting.GreetInput, *greeting.GreetOutput] {
	return temporalnexus.NewWorkflowRunOperation(greetingnexus.GreetingSlothGreetOperationName, server.SlothGreetWorkflow, func(ctx context.Context, input *greeting.GreetInput, options nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
		return client.StartWorkflowOptions{
			// Workflow IDs should typically be business meaningful IDs and are used to dedupe workflow starts.
			// For this example, we're using the request ID allocated by Temporal when the caller workflow schedules
			// the operation, this ID is guaranteed to be stable across retries of this operation.
			ID: options.RequestID,
			// Task queue defaults to the task queue this operation is handled on.
		}, nil
	})
}

func main() {
	clientOptions, err := options.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, taskQueue, worker.Options{})
	service, err := greetingnexus.NewGreetingNexusService(&handler{})
	if err != nil {
		log.Fatal(err)
	}
	w.RegisterNexusService(service)
	w.RegisterWorkflow(server.SlothGreetWorkflow)
	w.RegisterWorkflow(server.SlothSleepAndGreetWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
