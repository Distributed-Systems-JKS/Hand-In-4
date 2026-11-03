package main

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
	client := proto.NewClientServerClient(conn)

	// Context and request payload idk if we need this?
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()


	// Call remote requestAccess method
	res, err := client.RequestAccess(ctx, &proto.Request{Id: 0, Timestamp: 1})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	var permission string = "Access denied";
	if (res.GetPermission()) {
		permission = "Access granted"
	}

	log.Printf("Response from server: " + permission)
	
}


