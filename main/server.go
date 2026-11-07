package main

import (
	proto "Hand-In-4/grpc"
	"fmt"
	//"sync"
	//"config"
	"context"
	//"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	
)


type state int //state of the node

/*const ( //states: RELEASED, WANTED, HELD
	RELEASED state = iota //Not requesting the critical section
	WANTED                //Requesting the critical section, and waiting for permission from other nodes
	HELD                  //In the critical section
)*/

// Defining a server struct that embeds the unimplemented server
type server struct {
	proto.UnimplementedClientServerServer

}

// Implementing RequestAccess method defined in the .proto file
func (s *server) RequestAccess(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	node.mu.Lock()
	//Lamport clock sync
	node.lamportClock = max(node.lamportClock, req.GetTimestamp()) + 1
	algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] received access request from Node" + strconv.Itoa(int(req.Id)) + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
	
	//making a copy of the queue to ensure that it will not be changed by others
	ch := node.requestQueue

	//Defer the reply; if we are in HELD (CS) or WANTED (requesting CS), we will not reply to the request immediately, but instead defer the reply until we are done with the CS.
	useDefer := false

	if node.state == "HELD" {
		useDefer = true
		algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] is in state " + node.state + " and is queueing request from Node" + strconv.Itoa(int(req.Id)) + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
	} else if node.state == "WANTED" {
		//Check if the node in CS has a higher timestamp than the request we are currently processing. If it does, we will defer the reply until we are done with the CS.
		if node.timestamp < req.GetTimestamp() || (node.timestamp == req.GetTimestamp() && node.id < req.GetId()) {
			useDefer = true //defer if s.state == WANTED
			algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] is in state " + node.state + " request timestamp is " + strconv.Itoa(int(req.Timestamp)) + " and Node timestamp is " +  strconv.Itoa(int(node.timestamp)) + " so is queueing request from Node" + strconv.Itoa(int(req.Id)) + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
		} else{
			node.lamportClock++
			algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] is in state " + node.state + " request timestamp is " + strconv.Itoa(int(req.Timestamp)) + " and Node timestamp is " +  strconv.Itoa(int(node.timestamp)) + " so is replying to Node" + strconv.Itoa(int(req.Id)) + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
			return &proto.Response{Id: node.id, Timestamp: node.lamportClock, Permission: true}, nil // Permission her kommer an på response??
		}
	}


	node.mu.Unlock()
	if useDefer {
		<- ch //wait on the channel closing, when it is, then the code will continue running
	}
	

	node.lamportClock++
	algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] is in state " + node.state + " and is now replying to Node" + strconv.Itoa(int(req.Id)) + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
	return &proto.Response{Id: node.id, Timestamp: node.lamportClock, Permission: true}, nil 
}




func makeServer(port string) {
	
	// Listen on TCP port
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}


	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterClientServerServer(grpcServer, &server{}) 

	

	//algorithmLogger.Printf("Server listening on %s", port)

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
