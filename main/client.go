package main

import (
	proto "Hand-In-4/grpc"
	"strconv"

	"fmt"
	// "sync"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"google.golang.org/grpc/metadata"
)


func makeClient(port string) {

	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := proto.NewClientServerClient(conn)
	
	

	var clients[3] proto.ClientServerClient
	for i := 0; i < 3; i++ {
		var basePort string = "localhost:505" + strconv.Itoa(i);
		if(basePort != port){
			conn, err := grpc.NewClient(basePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
			log.Printf("Connecting to gRPC server on port 505" + strconv.Itoa(i))
			if err != nil {
				log.Fatalf("Did not connect: %v", err)
			}
			defer conn.Close()
			clients[i] = proto.NewClientServerClient(conn)
		} else {
			//not sure if this is used later on
			clients[i] = client;
		}
		
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
	//res, err := client.RequestAccess(context.Background(), &proto.Request{Id: 0, Timestamp: 1})
	if err != nil {
	//	log.Fatalf("Could not greet: %v", err)
	}
	for i:= 0; i < 3; i++{
		res, err := clients[i].RequestAccess(context.Background(),&proto.Request{Id: 0, Timestamp: 1}, grpc.WaitForReady(true))
		fmt.Println(res, err)
	}
	fmt.Println("ACCESS?");

	//log.Printf("Response from server: " + permission)

	
}


