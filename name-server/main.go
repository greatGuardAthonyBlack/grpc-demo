package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	option "grpc/echo-server/server"
	name_server "grpc/name-server/proto"
	"grpc/name-server/server"
	"log"
	"net"
	"time"
)

var (
	port = flag.Int("port", 50051, "server port")
)

func init() {
	flag.Parse()
}

func TestServer() {
	server.Register("echo", "localhost:50051")
	server.Register("echo", "localhost:50052")
	time.Sleep(3 * time.Second)
	server.Register("echo", "localhost:50053")
	time.Sleep(3 * time.Second)
	server.Register("echo", "localhost:50054")
	fmt.Println(server.GetNameServerCache())
}

func main() {
	//TestServer()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(option.GetServerOption()...)
	name_server.RegisterNameServerServer(s, &server.NameServer{})
	log.Printf("name server start listening at : %d", *port)
	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}

}
