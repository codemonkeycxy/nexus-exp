package main

import (
	"log"
	"os"

	tmplClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"nexus-exp/client"
	"nexus-exp/options"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	clientOptions, err := options.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := tmplClient.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, client.TaskQueue, worker.Options{})

	w.RegisterWorkflow(client.GreetWorkflow)
	w.RegisterWorkflow(client.SlothGreetWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
