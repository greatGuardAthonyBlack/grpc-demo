package client

import (
	"google.golang.org/grpc"
)

func GetClientOption() []grpc.DialOption {
	options := make([]grpc.DialOption, 0)

	options = append(options, GetTlsSecurityOption())
	return options

}
