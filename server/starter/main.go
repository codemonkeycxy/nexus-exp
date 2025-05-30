package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.temporal.io/api/enums/v1"
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
}

func (h *handler) Greet(name string) nexus.Operation[*greeting.GreetInput, *greeting.GreetOutput] {
	return nexus.NewSyncOperation(name, func(ctx context.Context, input *greeting.GreetInput, options nexus.StartOperationOptions) (*greeting.GreetOutput, error) {
		return &greeting.GreetOutput{
			Greeting: "Hello, " + input.Name,
		}, nil
	})
}

func (h *handler) SlothGreet(name string) nexus.Operation[*greeting.SlothGreetInput, *greeting.SlothGreetOutput] {
	return temporalnexus.MustNewWorkflowRunOperationWithOptions(temporalnexus.WorkflowRunOperationOptions[*greeting.SlothGreetInput, *greeting.SlothGreetOutput]{
		Name: name,
		Handler: func(ctx context.Context, input *greeting.SlothGreetInput, operationOptions nexus.StartOperationOptions) (temporalnexus.WorkflowHandle[*greeting.SlothGreetOutput], error) {
			return temporalnexus.ExecuteWorkflow(ctx, operationOptions, client.StartWorkflowOptions{
				ID:                       fmt.Sprintf("greet-sloth-%s", input.GetSlothName()),
				WorkflowIDConflictPolicy: enums.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING, // let the same sloth handle all greetings
			}, server.SlothSleepAndGreetWorkflow, server.SlothSleepAndGreetWorkflowInput{
				GreetInput: input,
				CountDown:  5,
			})
		},
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
	w.RegisterWorkflow(server.SlothSleepAndGreetWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
