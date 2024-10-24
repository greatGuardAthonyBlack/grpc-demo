package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	nc "grpc/name-client"
	"log"
)

var INSTANCES = []string{"localhost:50051", "localhost:50052", "localhost:50053"}
var nameServerClient *nc.NameClient
var nameServerAddr = "localhost:60051"

const (
	SCHEME  = "grpc"
	SERVICE = "echo"
)

type NameResolverBuilder struct {
}

func init() {

	nameServerClient = nc.BuildNameClient(nameServerAddr)
}

func (*NameResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := NameResolver{
		target:    target,
		cc:        cc,
		addrStore: map[string][]string{SERVICE: nameServerClient.GetServiceAddr(SERVICE)},
	}
	r.Start()
	return &r, nil
}

func (*NameResolverBuilder) Scheme() string {
	return SCHEME
}

func GetNameResolver() grpc.DialOption {
	return grpc.WithResolvers(&NameResolverBuilder{})
}

type NameResolver struct {
	target    resolver.Target
	cc        resolver.ClientConn
	addrStore map[string][]string
}

func (r *NameResolver) ResolveNow(opts resolver.ResolveNowOptions) {
	r.addrStore = map[string][]string{SERVICE: nameServerClient.GetServiceAddr(SERVICE)}
	log.Println("refresh client resolver")
	r.Start()
}

func (r *NameResolver) Close() {
}

func (r *NameResolver) Start() {
	addrs := r.addrStore[r.target.Endpoint()]
	addr_list := make([]resolver.Address, len(addrs))
	for _, el := range addrs {
		addr_list = append(addr_list, resolver.Address{
			Addr: el,
		})
	}
	r.cc.UpdateState(resolver.State{Addresses: addr_list})
}
