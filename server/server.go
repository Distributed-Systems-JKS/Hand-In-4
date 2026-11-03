//states: RELEASED, WANTED, HELD

package main;

import (
	proto "Hand-In-4/grpc"
	// "fmt"
	// "sync"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

)
// Defining a server struct that embeds the unimplemented server
type server struct {
	proto.UnimplementedGreeterServer
	// slice eller que med clients?

}

// Implementing SayHello method defined in the .proto file
func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error){
	log.Printf("Received request from: %s", req.GetName())
	return &proto.HelloResponse{Message: "Hello " + req.GetName() + "!"}, nil
}

func main() {
	// Listen on TCP port
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterGreeterServer(grpcServer, &server{})


	log.Printf("Server listening at %v", listener.Addr())

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}


