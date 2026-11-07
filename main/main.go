package main

import (
	//proto "Hand-In-4/grpc"
	//"fmt"
	"log"
	"os"

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
var algorithmLogger log.Logger

// Defining a server struct that embeds the unimplemented server
type Node struct {
	id           int64
	lamportClock int64      //changes for every local event
	mu           sync.Mutex //to protect state, timestamp, and lamportClock
	state        string      //which state is the node in
	timestamp    int64      //timestamp of the CURRENT request for the critical section, so important for state == WANTED
	requestQueue chan struct{}    //queue containing client ids to hold FIFO : keep track of requests from other nodes, so we can reply to them when current node is leaving CS

}


func makeLogger(){
	// --- Opens or creates a new server log
	const append = os.O_APPEND; // adss new logs to the end of the file
	const create = os.O_CREATE; // creates the file if it doesn't exist
	const wronly = os.O_WRONLY; // opens the file for writing only
	
	f, err := os.OpenFile("algorithm.log", append|create|wronly, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err) 
	}
	
	// New creates a new Logger with log data destination, prefix, logging properties.
	algorithmLogger = *log.New(f,"[LOG]", log.Ldate|log.Ltime)
}



func main(){
	makeLogger()
	portFlag := flag.String("port", "5050", "The port for this server to listen on")
	flag.Parse()
	port := *portFlag
	


	
	idInt, err := strconv.ParseInt(port[len(port)-1:], 10, 64) //get the last digit of the port and set as id
	if(err != nil){
		log.Printf("Could not convert id %v", err)
	}

	

	// -- INITIALISATION --
	node = Node{
		id: idInt,
		lamportClock: 0,
		state: "RELEASED",
		requestQueue: make(chan struct{}),
	}

	go makeServer(*portFlag)

	go makeClient(*portFlag)

	select{} //so our main continues to run even if the goroutine finishes temporarily

}