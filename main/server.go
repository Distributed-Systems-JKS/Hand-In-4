package main

import (
	proto "Hand-In-4/grpc"
	"fmt"

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

type state int //state of the node

const ( //states: RELEASED, WANTED, HELD
	RELEASED state = iota //Not requesting the critical section
	WANTED                //Requesting the critical section, and waiting for permission from other nodes
	HELD                  //In the critical section
)

//From previous assignment
/*
header := metadata.New(map[string]string{"client-id":strconv.Itoa(clientId)})
	grpc.SendHeader(header)
*/

// Defining a server struct that embeds the unimplemented server
type server struct {
	proto.UnimplementedClientServerServer

	id           int64
	lamportClock int64      //changes for every local event
	mu           sync.Mutex //to protect state, timestamp, and lamportClock
	state        state      //which state is the node in
	timestamp    int64      //timestamp of the CURRENT request for the critical section, so important for state == WANTED
	requestQueue []int64    //queue containing client ids to hold FIFO : keep track of requests from other nodes, so we can reply to them when current node is leaving CS

}

// Implementing RequestAccess method defined in the .proto file
func (s *server) RequestAccess(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	fmt.Println("Request received")
	s.mu.Lock()

	//Lamport clock update
	if req.GetTimestamp() > s.timestamp {
		s.timestamp = req.GetTimestamp()
	}
	s.timestamp++ //incrementing the timestamp of the current node, so it is always greater than the timestamp of the request
	//|| (req.GetTimestamp() == s.timestamp && req.GetId() > s.id) //if the timestamps are equal, we compare the ids of the nodes to break the tie

	//Defer the reply; if we are in HELD (CS) or WANTED (requesting CS), we will not reply to the request immediately, but instead defer the reply until we are done with the CS.
	//This is done by adding the request to a queue, and replying to it when we are done with the CS.
	useDefer := false

	if s.state == HELD {
		useDefer = true
	} else if s.state == WANTED {
		//Check if the node in CS has a higher timestamp than the request we are currently processing. If it does, we will defer the reply until we are done with the CS.
		if s.timestamp < req.GetTimestamp() || (s.timestamp == req.GetTimestamp() && s.id < req.GetId()) {
			useDefer = true //defer if s.state == WANTED
		}
	}

	if useDefer {
		//Save the id of the request in a queue, so we can reply to it later when we are done with the CS.
		s.requestQueue = append(s.requestQueue, req.GetId())
		s.mu.Unlock() //Unlocking the mutex before returning, so we can reply to the request later when we are done with the CS.
		log.Printf("Deferring reply to request from: %d", req.GetId())
		//Returning a response with permission set to false -> request is deferred and will be replied to later.
		return &proto.Response{Id: req.GetId(), Timestamp: req.GetTimestamp(), Permission: false}, nil // Permission her kommer an på response
	}

	s.mu.Unlock()

	log.Printf("Received request from: %d", req.GetId())
	return &proto.Response{Id: req.GetId(), Timestamp: req.GetTimestamp(), Permission: true}, nil // Permission her kommer an på response??
}

func (s *server) ReplyToRequest(ctx context.Context, res *proto.Response) (*proto.Empty, error) {
	s.mu.Lock()
	if res.GetTimestamp() > s.lamportClock {
		s.lamportClock = res.GetTimestamp() //updating the lamport clock of the current node
	}
	s.lamportClock++ //incrementing the lamport clock of the current node, so it is always greater than the timestamp of the request
	s.mu.Unlock()

	/*
		if s.state == WANTED {
			s.replyCount++ //incrementing the number of replies we have received from other nodes
			log.Printf("Received reply from: %d, total replies: %d", res.GetId(), s.replyCount)

			//If we are in WANTED state, we will check if we have received permission from all other nodes. If we have, we can enter the critical section.
			//We will check if we have received permission from all other nodes by checking if the number of replies we have received is equal to the number of nodes in the system.

			//Here we have to signal that we have received permission from all other nodes, so we can enter the critical section. This can be done by using a channel or a condition variable. For now, we will just log that we have received permission from all other nodes.
			//If we have received permission from all other nodes, we can enter the critical section.
			//We will also check if we have received permission from all other nodes by checking if the number of replies we have received is equal to the number of nodes in the system.
			}
	*/

	log.Printf("Received response from: %d", res.GetId())
	// return &proto.Response{Id: res.GetId(), Timestamp: res.GetTimestamp(), Permission: res.GetPermission()}, nil
	return &proto.Empty{}, nil
}

func makeServer(port string) {

	// Listen on TCP port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", port, err)
	}

	/*// Wait for new connections
	for {
		// Accept new connections
		c, err := listener.Accept()
		if err != nil {
			log.Printf("Listener returned: %s", err)
			break
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
		}()
	}*/

	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterClientServerServer(grpcServer, &server{}) //Its empty??

	log.Printf("Server listening on %s", port)

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
