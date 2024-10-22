package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc/echo"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	SESSION_TIMEOUT = 30
)

type EchoClient struct {
	conn *grpc.ClientConn
}

func BuildEchoClient(conn *grpc.ClientConn) *EchoClient {
	return &EchoClient{
		conn: conn,
	}
}

func GetClientDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), SESSION_TIMEOUT*time.Second)
}

func (c *EchoClient) CallUnary() {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}

	}()
	cli := echo.NewEchoClient(c.conn)
	ctx, cancel := GetClientDefaultContext()
	defer cancel()
	ctx = GetMetaContext(ctx, "el", "v1")
	req := &echo.EchoRequest{
		Message: "Unary request",
		Time:    timestamppb.New(time.Now()),
	}

	enhancedUnaryEcho := GetMetaDataMiddleware(cli.UnaryEcho)
	unaryEcho, err := enhancedUnaryEcho(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("unary echo:%s\n", unaryEcho.Message)

}

func (c *EchoClient) CallServerEchoStream() {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}

	}()
	cli := echo.NewEchoClient(c.conn)
	ctx, cancel := GetClientDefaultContext()
	defer cancel()
	req := &echo.EchoRequest{
		Message: "server stream request",
		Time:    timestamppb.New(time.Now()),
	}

	ctx = GetMetaContext(ctx, "el", "v1")
	stream, err := cli.ServerStreamEcho(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	savePath := "D:\\programing\\go_workspace\\grpc\\echo-client\\download\\" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
	file, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		file.Write(recv.Data[:recv.Len])

	}
	stream.CloseSend()
}

func (c *EchoClient) CallClientEchoStream() {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}

	}()
	path := "D:\\programing\\go_workspace\\grpc\\echo-client\\file\\client.jpg"
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cli := echo.NewEchoClient(c.conn)
	ctx, cancel := GetClientDefaultContext()
	defer cancel()
	ctx = GetMetaContext(ctx, "el", "v1")
	stream, err := cli.ClientStreamEcho(ctx)

	header, err := stream.Header()
	if err == nil {
		fmt.Printf("header metadata time:%s payload:%s\n", header.Get("time"), header.Get(HEADER_PAYLOAD_KEY))
	}

	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}

		request := &echo.EchoRequest{
			Data: buf[:n],
			Time: timestamppb.New(time.Now()),
			Len:  int32(n),
		}
		stream.Send(request)

	}

	recv, err := stream.CloseAndRecv()

	if err != nil {
		log.Println("client received failed", err)

	}
	tailer := stream.Trailer()
	fmt.Printf("header metadata time:%s payload:%s\n", header.Get("time"), tailer.Get(TAILER_PAYLOAD_KEY))
	fmt.Printf("client received end up data,payload message:%s\n", recv.Message)

}

func (c *EchoClient) CallWay2StreamEcho() {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("something happened :%s", err)
		}

	}()

	ctx, cancel := GetClientDefaultContext()
	defer cancel()
	ctx = GetMetaContext(ctx, "el", "v1")
	cli := echo.NewEchoClient(c.conn)
	stream, err := cli.Way2StreamEcho(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		path := "D:\\programing\\go_workspace\\grpc\\echo-client\\file\\client.jpg"
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				return
			}

			request := &echo.EchoRequest{
				Data: buf[:n],
				Time: timestamppb.New(time.Now()),
				Len:  int32(n),
			}
			stream.Send(request)

		}
		stream.CloseSend()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		savePath := "D:\\programing\\go_workspace\\grpc\\echo-client\\download\\" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
		file, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				return
			}
			file.Write(recv.Data[:recv.Len])

		}
	}()

	wg.Wait()

}
