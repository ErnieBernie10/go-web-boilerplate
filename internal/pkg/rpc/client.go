package rpc

import (
	pb "framer/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client *RpcClient

type RpcClient struct {
	Auth  pb.AppUserServiceClient
	Frame pb.FrameServiceClient
}

func NewClient(target string) error {
	conn, err := grpc.NewClient(target, grpc.DialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	if err != nil {
		return err
	}

	Client = &RpcClient{
		Auth:  pb.NewAppUserServiceClient(conn),
		Frame: pb.NewFrameServiceClient(conn),
	}

	return nil
}
