package main

import (
	//proto "Hand-In-4/grpc"
	"fmt"

	//"fmt"
	//"sync"
	//"config"
	//"context"
	"flag"
	//"log"
	//"net"
	//"strconv"

	//"google.golang.org/grpc"
	//"google.golang.org/grpc/metadata"
)

func main(){
	portFlag := flag.String("port", "5050", "The port for this server to listen on")
	flag.Parse()

	address := fmt.Sprintf(":%s", *portFlag)

	go makeServer(address)

	go makeClient(address)

	select{} //so our main continues to run even if the goroutine finishes temporarily

}