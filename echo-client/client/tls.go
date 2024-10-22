package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	ca_cert_path = "D:\\programing\\go_workspace\\grpc\\x509\\ca_cert.pem"
	domain       = "echo.grpc.0voice.com"
)

func GetDefaultSecurityOption() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func GetTlsSecurityOption() grpc.DialOption {
	cred, err := credentials.NewClientTLSFromFile(ca_cert_path, domain)
	if err != nil {
		log.Fatal(err)
	}
	option := grpc.WithTransportCredentials(cred)
	
	return option
}
