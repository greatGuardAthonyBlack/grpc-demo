package server

import (
	"google.golang.org/grpc/metadata"
	"time"
)

func getMetaByMap(mp map[string]string) metadata.MD {
	return metadata.New(mp)
}
func GetServerHeaderTailer() (header metadata.MD, tailer metadata.MD) {

	header = getMetaByMap(map[string]string{"time": time.Now().Format("2006-01-02T15:04Z07:00")})
	header.Set("h", "hv")

	tailer = getMetaByMap(map[string]string{"time": time.Now().Format("2006-01-02T15:04Z07:00")})
	tailer.Set("t", "tv")

	return header, tailer
}
