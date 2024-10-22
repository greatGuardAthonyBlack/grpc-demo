package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc/echo"
	"grpc/echo-server/server"
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
	s := grpc.NewServer(server.GetServerOption()...)
	echo.RegisterEchoServer(s, server.BuildEchoServer())
	log.Printf("echo server listening at : %d", *port)
	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}

}
