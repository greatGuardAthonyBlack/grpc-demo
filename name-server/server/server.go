package server

import (
	"context"
	"grpc/name"
	"io"
	"log"
)

type NameServer struct {
	name.UnimplementedNameServerServer
}

func (*NameServer) Register(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {

	for _, a := range in.Addr {
		Register(in.ServiceName, a)
		log.Printf("service: %s  addr:%s is registered\n", in.ServiceName, a)
	}

	return &name.NameResponse{
		ServiceName: in.ServiceName,
	}, nil
}

func (*NameServer) Delete(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	for _, a := range in.Addr {
		Delete(in.ServiceName, a)
	}

	return &name.NameResponse{
		ServiceName: in.ServiceName,
	}, nil
}
func (*NameServer) KeepAlive(stream name.NameServer_KeepAliveServer) error {
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for _, a := range recv.Addr {
			Keepalive(recv.ServiceName, a)
		}
	}
	return stream.SendAndClose(&name.NameResponse{})
}
func (*NameServer) GetAddr(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	list := make([]string, 0)
	list = append(list, GetService(in.ServiceName)...)
	return &name.NameResponse{
		ServiceName: in.ServiceName,
		Addr:        list,
	}, nil
}
