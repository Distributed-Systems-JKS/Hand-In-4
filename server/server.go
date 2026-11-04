//states: RELEASED, WANTED, HELD

package main

import (
	proto "Hand-In-4/grpc"
	"fmt"

	//"fmt"
	// "sync"
	//"config"
	"context"
	"log"
	"net"
	"flag"

	"google.golang.org/grpc"
)

/*
type state int //state of the node

const (
	RELEASED state = iota //Not requesting the critical section
	WANTED //Requesting the critical section, and waiting for permission from other nodes
	HELD //In the critical section
)
*/

// Defining a server struct that embeds the unimplemented server
type server struct {
	proto.UnimplementedClientServerServer

	/*
	id int64
	lamportClock int64
	mu sync.Mutex //to protect state, timestamp, and lamportClock
	state state //which state is the node in
	timestamp int64 //timestamp of the request for the critical section, so important for state == WANTED

	requestQueue []int 64 //queue containing client o to hold FIFO : keep track of requests from other nodes, so we can reply to them when current node is leaving CS
	
	?is there other fields missing??
	*/

}

// Implementing RequestAccess method defined in the .proto file
func (s *server) RequestAccess(ctx context.Context, req *proto.Request) (*proto.Response, error){
	log.Printf("Received request from: %d", req.GetId())
	return &proto.Response{Id: req.GetId(), Timestamp: req.GetTimestamp(), Permission: true}, nil // Permission her kommer an på response??
}

func (s *server) ReplyToRequest(ctx context.Context, res *proto.Response) (*proto.Empty, error){
	log.Printf("Received response from: %d", res.GetId())
	// return &proto.Response{Id: res.GetId(), Timestamp: res.GetTimestamp(), Permission: res.GetPermission()}, nil
	return &proto.Empty{}, nil
}


func main() {

	// Default port
	portFlag := flag.String("port", "5050", "The port for this server to listen on")
	flag.Parse()

	address := fmt.Sprintf(":%s", *portFlag)

	// Listen on TCP port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}

	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterClientServerServer(grpcServer, &server{})


	log.Printf("Server listening on %s", address)

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func startServer(address string){
	
}
