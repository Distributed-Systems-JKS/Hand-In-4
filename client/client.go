package main;

import (
	proto "Hand-In-4/grpc"
	// "fmt"
	// "sync"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func main() {

	// Establish connection
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create the gRPC client stub
	client := proto.NewGreeterClient(conn)

	// Context and request payload idk if we need this?
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call remote SayHello method
	res, err := client.SayHello(ctx, &proto.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	log.Printf("Response from server: %s", res.GetMessage())
	
}


