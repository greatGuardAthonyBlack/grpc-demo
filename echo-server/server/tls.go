package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

const server_cert_path = "D:\\programing\\go_workspace\\grpc\\x509\\server_cert.pem"
const server_cert_key_path = "D:\\programing\\go_workspace\\grpc\\x509\\server_key.pem"

func getTlsOption() grpc.ServerOption {
	cred, err := credentials.NewServerTLSFromFile(server_cert_path, server_cert_key_path)
	if err != nil {
		log.Fatal(err)
	}
	return grpc.Creds(cred)
}
