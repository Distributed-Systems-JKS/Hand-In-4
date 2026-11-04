package client

import (
	proto "Hand-In-4/grpc"
	//"strconv"

	// "fmt"
	// "sync"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"google.golang.org/grpc/metadata"
)



func RunClient(connections map[int]string, address string, myId int64) {
	
	clients := make(map[int]proto.ClientServerClient)

	for id, addr := range connections {

		log.Printf("Connecting to gRPC server on port" + address)

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		clients[id] = proto.NewClientServerClient(conn)

	}


	for targetId := range connections {
		// Skip own id
		
		log.Printf("Client %d sending RequestAccess to Node %d...", myId, targetId)
		res, err := clients[targetId].RequestAccess(context.Background(), &proto.Request{
			Id: 		myId, 
			Timestamp: 	1,
		})
		if err != nil {
			log.Fatalf("Could not request access: %v", err)
		}

		var permission string = "Access denied";
		if (res.GetPermission()) {
			permission = "Access granted"
		}

		log.Printf("Response from server: " + permission)
	


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

	