package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporalnexus"
	"go.temporal.io/sdk/workflow"

	greeting "nexus-exp/gen/proto/v1"
	"nexus-exp/gen/proto/v1/greetingnexus"

	"github.com/nexus-rpc/sdk-go/nexus"
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
	return temporalnexus.NewWorkflowRunOperation(greetingnexus.GreetingSlothGreetOperationName, HelloHandlerWorkflow, func(ctx context.Context, input *greeting.GreetInput, options nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
		return client.StartWorkflowOptions{
			// Workflow IDs should typically be business meaningful IDs and are used to dedupe workflow starts.
			// For this example, we're using the request ID allocated by Temporal when the caller workflow schedules
			// the operation, this ID is guaranteed to be stable across retries of this operation.
			ID: options.RequestID,
			// Task queue defaults to the task queue this operation is handled on.
		}, nil
	})
}

func HelloHandlerWorkflow(ctx workflow.Context, input *greeting.GreetInput) (*greeting.GreetOutput, error) {
	if err := workflow.Sleep(ctx, 5*time.Second); err != nil {
		return nil, err
	}

	return &greeting.GreetOutput{
		Greeting: "Hello, " + input.Name,
	}, nil
}

func main() {
	service, err := greetingnexus.NewGreetingNexusService(&handler{})
	if err != nil {
		log.Fatal(err)
	}
	registry := nexus.NewServiceRegistry()
	if err := registry.Register(service); err != nil {
		log.Fatal(err)
	}
	rh, err := registry.NewHandler()
	if err != nil {
		log.Fatal(err)
	}
	h := nexus.NewHTTPHandler(nexus.HandlerOptions{
		Handler:    rh,
		Serializer: nexus.DefaultSerializer(),
	})

	listener, err := net.Listen("tcp", "localhost:7243")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	if err = http.Serve(listener, h); err != nil {
		log.Fatal(err)
	}
}
