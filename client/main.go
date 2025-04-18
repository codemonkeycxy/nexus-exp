package main

import (
	"context"
	"log"
	"time"

	greeting "nexus-exp/gen/proto/v1"
	"nexus-exp/gen/proto/v1/greetingnexus"

	"github.com/nexus-rpc/sdk-go/nexus"
)

func main() {
	client, err := greetingnexus.NewGreetingNexusHTTPClient(nexus.HTTPClientOptions{
		BaseURL: "http://localhost:7243",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	result, err := client.Greet(ctx, &greeting.GreetInput{
		Name: "World",
	}, nexus.ExecuteOperationOptions{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received greeting: %s", result.Greeting)

	now := time.Now()
	result, err = client.SlothGreet(ctx, &greeting.GreetInput{
		Name: "World",
	}, nexus.ExecuteOperationOptions{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received greeting from sloth after %s: %s", time.Since(now), result.Greeting)
}
