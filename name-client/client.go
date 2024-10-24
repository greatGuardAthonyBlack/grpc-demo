package name_client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	clientPool "grpc/echo-client/pool"
	"grpc/name"
	"log"
	"time"
)

const (
	SESSION_TIMEOUT    = 30
	HEARTBEAT_INTERVAL = 15
	NAME_SERVICE_ADDR  = "localhost:60051"
)

var pool *clientPool.ClientPool

func GetClientPool() *clientPool.ClientPool {
	if pool == nil {
		buildPool, _ := clientPool.BuildPool(NAME_SERVICE_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
		pool = buildPool
	}
	return pool
}

func init() {
	buildPool, _ := clientPool.BuildPool(NAME_SERVICE_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	pool = buildPool
}

type NameClient struct {
	conn *grpc.ClientConn
}

func GetClientDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), SESSION_TIMEOUT*time.Second)
}

func BuildNameClient(addr string) *NameClient {
	pool := GetClientPool()
	return &NameClient{
		conn: pool.Get(),
	}
}

func (c *NameClient) Register(serviceName string, addr string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}
	}()
	ctx, cancelFunc := GetClientDefaultContext()
	defer cancelFunc()

	req := &name.NameRequest{
		ServiceName: serviceName,
		Addr:        []string{addr},
	}
	cli := name.NewNameServerClient(c.conn)
	_, err := cli.Register(ctx, req)
	if err != nil {
		log.Print("register service failed : ", err)
	}

}

func (c *NameClient) Delete(serviceName string, addr string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}
	}()
	ctx, cancelFunc := GetClientDefaultContext()
	defer cancelFunc()

	req := &name.NameRequest{
		ServiceName: serviceName,
		Addr:        []string{addr},
	}
	cli := name.NewNameServerClient(c.conn)
	_, err := cli.Delete(ctx, req)
	if err != nil {
		log.Print(err)
	}
}

func (c *NameClient) GetServiceAddr(serviceName string) []string {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}
	}()
	ctx, cancelFunc := GetClientDefaultContext()
	defer cancelFunc()

	req := &name.NameRequest{
		ServiceName: serviceName,
	}
	cli := name.NewNameServerClient(c.conn)
	resp, err := cli.GetAddr(ctx, req)
	if err != nil {
		log.Print(err)
		return []string{}
	}

	return resp.GetAddr()
}

func (c *NameClient) Keepalive(serviceName string, addr string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}
	}()

	req := &name.NameRequest{
		ServiceName: serviceName,
		Addr:        []string{addr},
	}
	cli := name.NewNameServerClient(c.conn)
	stream, err := cli.KeepAlive(context.Background())
	if err != nil {
		log.Print(err)
	}
	for {
		stream.Send(req)
		time.Sleep(HEARTBEAT_INTERVAL * time.Second)
	}

}
