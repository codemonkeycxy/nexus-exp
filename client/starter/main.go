package main

import (
	"context"
	"log"
	"os"
	"time"

	tmprlClient "go.temporal.io/sdk/client"

	"nexus-exp/client"
	"nexus-exp/options"
)

func main() {
	clientOptions, err := options.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := tmprlClient.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	runWorkflow(c, client.GreetWorkflow, "World")
	runWorkflow(c, client.SlothGreetWorkflow, "World", []string{"Snugglemuffin", "Snugglemuffin", "Snugglemuffin"}) // Greet the same sloth 3 times
	runWorkflow(c, client.SlothGreetWorkflow, "World", []string{"Snugglemuffin", "Mochapaws", "Lazeberry"})         // Greet 3 different sloths
}

func runWorkflow(c tmprlClient.Client, workflow interface{}, args ...interface{}) {
	ctx := context.Background()
	workflowOptions := tmprlClient.StartWorkflowOptions{
		ID:        "nexus_greet_caller_workflow_" + time.Now().Format("20060102150405"),
		TaskQueue: client.TaskQueue,
	}

	wr, err := c.ExecuteWorkflow(ctx, workflowOptions, workflow, args...)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = wr.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
