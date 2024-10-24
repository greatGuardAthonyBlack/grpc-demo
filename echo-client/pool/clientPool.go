package pool

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"sync"
)

type ClientPoolApi interface {
	Get() *grpc.ClientConn
	Put(*grpc.ClientConn)
}
type ClientPool struct {
	pool sync.Pool
}

func BuildPool(target string, opts ...grpc.DialOption) (*ClientPool, error) {
	return &ClientPool{
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

func (p *ClientPool) Get() *grpc.ClientConn {
	conn := p.pool.Get().(*grpc.ClientConn)
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = p.pool.New().(*grpc.ClientConn)
	}
	return conn
}

func (p *ClientPool) Put(conn *grpc.ClientConn) {
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = p.pool.New().(*grpc.ClientConn)
	}
	p.pool.Put(conn)
}
