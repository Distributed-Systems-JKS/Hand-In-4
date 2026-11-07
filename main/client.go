package main

import (
	proto "Hand-In-4/grpc"
	"strconv"
	"fmt"
	"sync"
	"context"
	"log"
	"time"
	"math/rand"
	

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"google.golang.org/grpc/metadata"
)



func makeClient(port string) {

	address := fmt.Sprintf("localhost:%s", port)
	
	//conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//client := proto.NewClientServerClient(conn)

	//making the other client stubs, so we can communicate to the other servers
	var clients[2] proto.ClientServerClient
	clientPosition := 0 //making sure that we only add two clients to the list in the for-loop
	for i := 0; i < 3; i++ {
		var basePort string = "localhost:505" + strconv.Itoa(i);
		if(basePort != address){
			conn, err := grpc.NewClient(basePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
			//log.Printf("Connecting to gRPC server on port 505" + strconv.Itoa(i))
			if err != nil {
				log.Fatalf("Did not connect: %v", err)
			}
			clients[clientPosition] = proto.NewClientServerClient(conn)
			clientPosition++
		}
		
	}


	go enter(clients)

	

	
}

func enter(clients[2] proto.ClientServerClient) {
	var wg sync.WaitGroup

	timer := rand.Intn(7) + 3 //random int from 3 to 10
	algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] waiting " + strconv.Itoa(int(timer)) + " until requesting access at lamport clock " + strconv.Itoa(int(node.lamportClock)))
	time.Sleep(time.Duration(timer)*time.Second)
	
	node.mu.Lock()
	node.state = "WANTED"
	node.timestamp = node.lamportClock
	node.requestQueue = make(chan struct{})

	// -- UPDATING LAMPORT CLOCK --
	node.lamportClock++;
	node.timestamp = node.lamportClock;
	node.mu.Unlock()


	algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] state is " + node.state + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))

	//adding 2 to the waitgroup, so the rest of the code finishes when 2 go-routines are done
	wg.Add(2)
	algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +"] is sending access request with timestamp " + strconv.Itoa(int(node.timestamp)) + " to all other nodes at lamport clock " + strconv.Itoa(int(node.lamportClock)))
	for i := range clients{
		go func(i int){
			defer wg.Done()
			reply, err := clients[i].RequestAccess(context.Background(), &proto.Request{Id: node.id, Timestamp: node.timestamp}, grpc.WaitForReady(true))
		
		if err != nil {
			log.Printf("Error with reply", err)
			return
		}

		node.mu.Lock()
		node.lamportClock = max(node.lamportClock,reply.GetTimestamp()) + 1
		node.mu.Unlock()

		algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] got a reply from Node" + strconv.Itoa(int(reply.Id)) + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
		time.Sleep(time.Duration(1)*time.Second)

	}(i)
	}
	
	wg.Wait() //Code after this can be run after the twi go-routines (the reply to request) in the loop are finished
	
	inCriticalSection()

	timer = rand.Intn(10)

	time.Sleep(time.Duration(timer)*time.Second)

	leaveCriticalSection()

	fmt.Println("Has been in Critical Section")
	
	
}


func inCriticalSection(){
	node.mu.Lock()
	node.state = "HELD"
	// -- UPDATING LAMPORT CLOCK --
	node.lamportClock++
	node.mu.Unlock()
	algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] is in Critical Section and state is " + node.state + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
}

func leaveCriticalSection(){
	node.mu.Lock()
	node.state = "RELEASED"
	// -- UPDATING LAMPORT CLOCK --
	node.lamportClock++
	node.mu.Unlock()
	algorithmLogger.Println("[Node" + strconv.Itoa(int(node.id)) +  "] left Critical Section and state is " + node.state + " at lamport clock " + strconv.Itoa(int(node.lamportClock)))
	node.mu.Lock()
	close(node.requestQueue)
	node.mu.Unlock()
}


