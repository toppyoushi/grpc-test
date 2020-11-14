package helloworld

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/toppyoushi/grpc-apis/pkg/helloworld"
)

type GreeterServerImp struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServerImp) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloRsp, error) {
	if in == nil {
		log.WithField("action", "receive request").Error("nil request")
		return nil, fmt.Errorf("invalid parameter")
	}

	rsp := &pb.HelloRsp{
		Msg: "Hello, I am a grpc server",
	}

	return rsp, nil
}
