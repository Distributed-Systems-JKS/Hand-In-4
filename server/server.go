//states: RELEASED, WANTED, HELD

package main;

import (
	proto "Hand-In-4/grpc"
	// "fmt"
	// "sync"
	"log"
	"net"
	"context"

	"google.golang.org/grpc"

)
// Defining a server struct that embeds the unimplemented server
type server struct {
	proto.UnimplementedClientServer
	// slice eller que med clients?

}

// Implementing RequestAccess method defined in the .proto file
func (s *server) RequestAccess(ctx context.Context, req *proto.Request) (*proto.Response, error){
	log.Printf("Received request from: %d", req.GetId())
	return &proto.Response{Id: req.GetId(), Timestamp: req.GetTimestamp(), Permission: true}, nil // Permission her kommer an på response??
}

func (s *server) ReplyToRequest(ctx context.Context, res *proto.Response) (*proto.Response, error){
	log.Printf("Received response from: %d", res.GetId())
	return &proto.Response{Id: res.GetId(), Timestamp: res.GetTimestamp(), Permission: res.GetPermission()}, nil
}



func main() {
	// Listen on TCP port
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterClientServer(grpcServer, &server{})


	log.Printf("Server listening at %v", listener.Addr())

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}


