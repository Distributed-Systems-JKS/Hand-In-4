package server

import (
	proto "Hand-In-4/grpc"
	//"fmt"

	//"fmt"
	"sync"
	//"config"
	"context"
	//"flag"
	"log"
	"net"
	//"strconv"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/metadata"
)

/*
type state int //state of the node

const ( //states: RELEASED, WANTED, HELD
	RELEASED state = iota //Not requesting the critical section
	WANTED //Requesting the critical section, and waiting for permission from other nodes
	HELD //In the critical section
)
*/


//From previous assignment
/*
header := metadata.New(map[string]string{"client-id":strconv.Itoa(clientId)})
	grpc.SendHeader(header) 
	*/


// Defining a server struct that embeds the unimplemented server
type server struct {
	proto.UnimplementedClientServerServer
	id int64
	lamportClock int64 //changes for every local event 
	mu sync.Mutex //to protect state, timestamp, and lamportClock
	//state state //which state is the node in
	timestamp int64 //timestamp of the CURRENT request for the critical section, so important for state == WANTED

	//requestQueue []int 64 //queue containing client o to hold FIFO : keep track of requests from other nodes, so we can reply to them when current node is leaving CS
	
	//?is there other fields missing??
	

}

// Implementing RequestAccess method defined in the .proto file
func (s *server) RequestAccess(ctx context.Context, req *proto.Request) (*proto.Response, error){
	s.mu.Lock()
	if (req.GetTimestamp() > s.timestamp){ 
		s.timestamp = req.GetTimestamp()
	}	
	s.timestamp++ //incrementing the timestamp of the current node, so it is always greater than the timestamp of the request
	//|| (req.GetTimestamp() == s.timestamp && req.GetId() > s.id) //if the timestamps are equal, we compare the ids of the nodes to break the tie

	//something something - K

	//compare the timestamps of the request and the current node, and if the request has a higher timestamp, we update the timestamp of the current node to be the same as the request's timestamp. This ensures that the current node's timestamp is always greater than or equal to any incoming request's timestamp.
	//defer if s.state == WANTED ??
	
	s.mu.Unlock()
	
	log.Printf("Received request from: %d", req.GetId())
	return &proto.Response{Id: req.GetId(), Timestamp: req.GetTimestamp(), Permission: true}, nil // Permission her kommer an på response??
}

func (s *server) ReplyToRequest(ctx context.Context, res *proto.Response) (*proto.Empty, error){
	s.mu.Lock() 
	if (res.GetTimestamp() > s.timestamp){
		s.timestamp = res.GetTimestamp()
	}
	s.timestamp++ //incrementing the timestamp since a message is always an event
	s.mu.Unlock()

	//something something - K
	
	log.Printf("Received response from: %d", res.GetId())
	// return &proto.Response{Id: res.GetId(), Timestamp: res.GetTimestamp(), Permission: res.GetPermission()}, nil
	return &proto.Empty{}, nil
}


func StartServer(address string, id int64) {

	
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterClientServerServer(grpcServer, &server{id: id})

	log.Printf("Server listening on %s", address)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}


}

/*
	// Wait for new connections
	for {
		// Accept new connections
		c, err := listener.Accept()
		if err != nil {
			log.Printf("Listener returned: %s", err)
			break
		}

		// Creating a new gRPC server instance
		




	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}



		
		// Start new Goroutine to handle new connection
		go func() {
			//defer c.Close()
			log.Printf("New connection created")

			// Writing something for now and closing the sessions
			_, err := c.Write([]byte("Helloooo"))
			if err != nil {
				log.Printf("Unable to write on connection: %s", err)
			}
		}()*/
