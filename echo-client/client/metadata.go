package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc/echo"
	"time"
)

const (
	HEADER_PAYLOAD_KEY = "h"
	TAILER_PAYLOAD_KEY = "t"
)

type ClientAPI func(ctx context.Context, in *echo.EchoRequest, opts ...grpc.CallOption) (*echo.EchoResponse, error)

func getMetaByMap(mp map[string]string) metadata.MD {
	return metadata.New(mp)
}

func getMetaByKV(kv ...string) metadata.MD {
	return metadata.Pairs(kv...)
}

func getOutcomingContext(ctx context.Context, mp map[string]string) context.Context {
	return metadata.NewOutgoingContext(ctx, getMetaByMap(mp))
}

func AppendKV2MetaContext(ctx context.Context, kv ...string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, kv...)
}

func GetMetaContext(ctx context.Context, kv ...string) context.Context {
	mp := map[string]string{"time": time.Now().Format("2006-01-02T15:04Z07:00")}
	ctx = getOutcomingContext(ctx, mp)
	if kv != nil {
		ctx = AppendKV2MetaContext(ctx, kv...)
	}
	return ctx

}

func GetMetaDataMiddleware(clientApi ClientAPI) ClientAPI {
	return func(ctx context.Context, in *echo.EchoRequest, opts ...grpc.CallOption) (*echo.EchoResponse, error) {
		var header, tailer metadata.MD
		opts = append(opts, grpc.Header(&header), grpc.Trailer(&tailer))
		res, err := clientApi(ctx, in, opts...)
		fmt.Printf("header metadata time:%s payload:%s\n", header.Get("time"), header.Get(HEADER_PAYLOAD_KEY))
		fmt.Printf("tailer metadata time:%s payload:%s\n", header.Get("time"), tailer.Get(TAILER_PAYLOAD_KEY))
		return res, err
	}
}
