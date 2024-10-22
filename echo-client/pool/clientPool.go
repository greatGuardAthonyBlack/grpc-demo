package pool

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"sync"
)

type ClientPool interface {
	Get() *grpc.ClientConn
	Put(*grpc.ClientConn)
}
type clientPool struct {
	pool sync.Pool
}

func BuildPool(target string, opts ...grpc.DialOption) (*clientPool, error) {
	return &clientPool{
		pool: sync.Pool{
			New: func() any {
				conn, err := grpc.NewClient(target, opts...)
				if err != nil {
					log.Print(err)
					return nil
				}
				return conn
			},
		},
	}, nil
}

func (p *clientPool) Get() *grpc.ClientConn {
	conn := p.pool.Get().(*grpc.ClientConn)
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = p.pool.New().(*grpc.ClientConn)
	}
	return conn
}

func (p *clientPool) Put(conn *grpc.ClientConn) {
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = p.pool.New().(*grpc.ClientConn)
	}
	p.pool.Put(conn)
}
