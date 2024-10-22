package main

import (
	"flag"
	"grpc/echo-client/client"
	pool2 "grpc/echo-client/pool"
	"log"
)

var (
	addr = flag.String("addr", "localhost:50051", "target server ")
)

func init() {
	flag.Parse()

}

func main() {

	connectPool, err := pool2.BuildPool(*addr, client.GetClientOption()...)
	if err != nil {
		log.Fatal(err)
	}

	conn := connectPool.Get()
	defer connectPool.Put(conn)

	echoClient := client.BuildEchoClient(conn)
	echoClient.CallUnary()
	echoClient.CallClientEchoStream()
	//echoClient.CallClientEchoStream()
	echoClient.CallWay2StreamEcho()
}
