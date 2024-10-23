package server

import (
	"context"
	name_server "grpc/name-server/proto"
	"io"
)

type NameServer struct {
	name_server.UnimplementedNameServerServer
}

func (*NameServer) Register(ctx context.Context, in *name_server.NameRequest) (*name_server.NameResponse, error) {

	for _, a := range in.Addr {
		Register(in.ServiceName, a)
	}

	return &name_server.NameResponse{
		ServiceName: in.ServiceName,
	}, nil
}

func (*NameServer) Delete(ctx context.Context, in *name_server.NameRequest) (*name_server.NameResponse, error) {
	for _, a := range in.Addr {
		Delete(in.ServiceName, a)
	}

	return &name_server.NameResponse{
		ServiceName: in.ServiceName,
	}, nil
}
func (*NameServer) KeepAlive(stream name_server.NameServer_KeepAliveServer) error {
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
	return stream.SendAndClose(&name_server.NameResponse{})
}
func (*NameServer) GetAddr(ctx context.Context, in *name_server.NameRequest) (*name_server.NameResponse, error) {
	list := make([]string, 0)
	list = append(list, GetService(in.ServiceName)...)
	return &name_server.NameResponse{
		ServiceName: in.ServiceName,
		Addr:        list,
	}, nil
}
