package main

import (
	//proto "Hand-In-4/grpc"
	"fmt"
	"log"

	//"fmt"
	//"sync"
	//"config"
	//"context"
	"flag"
	//"log"
	//"net"
	"strconv"
	"sync"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/metadata"
)

var node Node

// Defining a server struct that embeds the unimplemented server
type Node struct {

	id           int64
	lamportClock int64      //changes for every local event
	mu           sync.Mutex //to protect state, timestamp, and lamportClock
	state        state      //which state is the node in
	timestamp    int64      //timestamp of the CURRENT request for the critical section, so important for state == WANTED
	requestQueue []int64    //queue containing client ids to hold FIFO : keep track of requests from other nodes, so we can reply to them when current node is leaving CS

}



func main(){
	portFlag := flag.String("port", "5050", "The port for this server to listen on")
	flag.Parse()

	address := fmt.Sprintf(":%s", *portFlag)
	
	idInt, err := strconv.ParseInt(address[len(address)-1:], 10, 64) //get the last digit of the port and set as id
	if(err != nil){
		log.Printf("Could not convert id %v", err)
	}

	// -- INITIALISATION --
	node := &Node{
		id: idInt,
		lamportClock: 0,
		state: RELEASED,
	}

	go makeServer(address)

	go makeClient(address)

	select{} //so our main continues to run even if the goroutine finishes temporarily

}