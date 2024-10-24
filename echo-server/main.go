package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc/echo"
	"grpc/echo-server/server"
	name_client "grpc/name-client"
	"log"
	"net"
	"os"
	"os/signal"
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

	//start echo server
	go func() {
		if err = s.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	registerAddr := fmt.Sprintf("localhost:%d", *port)
	//keepalive
	nameClient := name_client.BuildNameClient("localhost:60051")
	go func() {
		nameClient.Register("echo", registerAddr)
		nameClient.Keepalive("echo", registerAddr)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()
	<-ctx.Done()

}
