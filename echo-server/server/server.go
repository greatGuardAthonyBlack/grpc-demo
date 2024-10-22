package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	CONTEXT_PAYLOAD_KEY = "el"
)

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func HandleContextMeta(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	fmt.Printf("request meta data time:%s   payload:%s\n", md.Get("time"), md.Get(CONTEXT_PAYLOAD_KEY))

}
func BuildEchoServer() EchoServer {
	return EchoServer{}
}
func (EchoServer) UnaryEcho(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	fmt.Printf("UnaryEcho api received: %s\n", in.Message)
	HandleContextMeta(ctx)

	header, tailer := GetServerHeaderTailer()
	defer func() {
		grpc.SetTrailer(ctx, tailer)
	}()
	grpc.SendHeader(ctx, header)
	resp := &echo.EchoResponse{
		Message: fmt.Sprintf("echo for message: %s", in.Message),
		Time:    timestamppb.New(time.Now()),
	}
	return resp, nil
}
func (EchoServer) ServerStreamEcho(in *echo.EchoRequest, stream echo.Echo_ServerStreamEchoServer) error {
	fmt.Printf("ServerStreamEcho api received: %s\n", in.Message)
	HandleContextMeta(stream.Context())

	header, tailer := GetServerHeaderTailer()
	defer func() {
		stream.SetTrailer(tailer)
	}()
	stream.SendHeader(header)
	resourcePath := "D:\\programing\\go_workspace\\grpc\\echo-server\\file\\server.jpg"
	f, err := os.Open(resourcePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		resp := &echo.EchoResponse{
			Data: buf[:n],
			Time: timestamppb.New(time.Now()),
			Len:  int32(n),
		}
		stream.Send(resp)
	}
	return nil
}
func (EchoServer) ClientStreamEcho(stream echo.Echo_ClientStreamEchoServer) error {
	savePath := "D:\\programing\\go_workspace\\grpc\\echo-server\\upload\\" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
	f, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	HandleContextMeta(stream.Context())

	header, tailer := GetServerHeaderTailer()
	defer func() {
		stream.SetTrailer(tailer)
	}()
	stream.SendHeader(header)

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print("stream received failed:", err)
			return err
		}
		f.Write(recv.Data[:recv.Len])

	}
	endResponse := &echo.EchoResponse{
		Message: "stream received success",
		Time:    timestamppb.New(time.Now()),
	}

	return stream.SendAndClose(endResponse)
}
func (EchoServer) Way2StreamEcho(stream echo.Echo_Way2StreamEchoServer) error {
	HandleContextMeta(stream.Context())

	header, tailer := GetServerHeaderTailer()
	defer func() {
		stream.SetTrailer(tailer)
	}()
	stream.SendHeader(header)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		savePath := "D:\\programing\\go_workspace\\grpc\\echo-server\\upload\\" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
		f, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Print("stream received failed:", err)
				return
			}
			f.Write(recv.Data[:recv.Len])

		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		resourcePath := "D:\\programing\\go_workspace\\grpc\\echo-server\\file\\server.jpg"
		f, err := os.Open(resourcePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
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
			resp := &echo.EchoResponse{
				Data: buf[:n],
				Time: timestamppb.New(time.Now()),
				Len:  int32(n),
			}
			stream.Send(resp)
		}
	}()
	wg.Wait()
	return nil
}
