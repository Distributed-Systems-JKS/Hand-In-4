package main

import (
	proto "Hand-In-4/grpc"
	"Hand-In-4/client"
	"Hand-In-4/server"

	//"log"
	//"strconv"
	//"google.golang.org/grpc"


)

var clients[3] proto.ClientServerClient


func main (){

	nodes := map[int]string {
		0: "localhost:5050",
		1: "localhost:5051",
		2: "localhost:5052",
	}

	for i := 0; i < 3; i++ {
		address := nodes[i]
		id := int64(i)

		go server.StartServer(address, id)

		// Give server a moment?? with time.Sleep

	}

	for i := 0; i < 3; i++ {
		address := nodes[i]
		id := int64(i)

		client.RunClient(nodes, address, id)


	}

	select {}


}

