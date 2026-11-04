package main

import (
	proto "Hand-In-4/grpc"
	"strconv"

	// "fmt"
	// "sync"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"google.golang.org/grpc/metadata"
)


func main() {

	// Default address and port 
	var basePort string = "localhost:505";

	var clients[3] proto.ClientServerClient
	for i := 0; i < 3; i++ {
		log.Printf("Connecting to gRPC server on port 505" + strconv.Itoa(i))
		conn, err := grpc.NewClient(basePort + strconv.Itoa(i), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()
		clients[i] = proto.NewClientServerClient(conn)
	}

	

//From previous assignment 
/*Doesn't work the same way as before, since we are not using a stream anymore.
	header, err := streamMessage.Header()
	if err != nil {
		log.Fatalf("failed to get header from stream: %v", err) 
	}
	id := header["client-id"][0] 
	*/

	


	// Create the gRPC client stub
	// client := proto.NewClientServerClient(conn)

	
	// Call remote requestAccess method
	res, err := clients[0].RequestAccess(context.Background(), &proto.Request{Id: 0, Timestamp: 1})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	var permission string = "Access denied";
	if (res.GetPermission()) {
		permission = "Access granted"
	}

	log.Printf("Response from server: " + permission)
	
}


