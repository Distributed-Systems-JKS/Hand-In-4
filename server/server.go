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
	proto.UnimplementedClientServerServer
	// slice eller que med clients?

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
	
	// Listen on Node 1 port
	listen1, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Listen on Node 2 port
	listen2, err := net.Listen("tcp", ":5052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Listen on Node 3 port
	listen3, err := net.Listen("tcp", ":5053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}




	// Creating a new gRPC server instance for Node 1
	grpcServer1 := grpc.NewServer()
	proto.RegisterClientServerServer(grpcServer1, &server{})
	log.Printf("Server listening at Node 1: %v", listen1.Addr())

	// Creating a new gRPC server instance for Node 2
	grpcServer2 := grpc.NewServer()
	proto.RegisterClientServerServer(grpcServer2, &server{})
	log.Printf("Server listening at Node 2: %v", listen2.Addr())

	// Creating a new gRPC server instance for Node 3
	grpcServer3 := grpc.NewServer()
	proto.RegisterClientServerServer(grpcServer3, &server{})
	log.Printf("Server listening at Node 3: %v", listen3.Addr())





	// Start serving requests for Node 1
	if err := grpcServer1.Serve(listen1); err != nil {
		log.Fatalf("Failed to serve Node 1: %v", err)
	}

	// Start serving requests for Node 2
	if err := grpcServer2.Serve(listen2); err != nil {
		log.Fatalf("Failed to serve Node 2: %v", err)
	}

	// Start serving requests for Node 3
	if err := grpcServer3.Serve(listen3); err != nil {
		log.Fatalf("Failed to serve Node 3: %v", err)
	}

}


