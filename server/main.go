package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

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
	return nexus.NewSyncOperation(name, func(ctx context.Context, input *greeting.GreetInput, options nexus.StartOperationOptions) (*greeting.GreetOutput, error) {
		time.Sleep(1 * time.Second)
		return &greeting.GreetOutput{
			Greeting: "Hello, " + input.Name,
		}, nil
	})
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
