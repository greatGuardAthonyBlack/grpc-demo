package server

import "google.golang.org/grpc"

func GetServerOption() []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0)
	opts = append(opts, getTlsOption())
	return opts
}
