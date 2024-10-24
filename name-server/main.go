package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	name_server "grpc/name"
	"grpc/name-server/server"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "server port")
)

func init() {
	flag.Parse()
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	name_server.RegisterNameServerServer(s, &server.NameServer{})
	log.Printf("name server start listening at : %d", *port)
	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}

}
