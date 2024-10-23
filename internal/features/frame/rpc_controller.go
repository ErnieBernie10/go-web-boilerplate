package frame

import (
	"context"
	"framer/internal/database"
	pb "framer/internal/proto"
	"framer/internal/rpc"

	"github.com/google/uuid"
)

type FrameController struct {
	pb.FrameServiceServer
}

func (c *FrameController) CreateFrame(ctx context.Context, req *pb.CreateFrameRequest) (*pb.Frame, error) {
	claims, err := rpc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	frame, err := fromRpc(req, claims.ID, uuid.NullUUID{})
	if err != nil {
		return nil, err
	}

	id, err := SaveFrameWithFile(ctx, database.Service, frame, req.Frame.File.Content, req.Frame.File.FileName)
	if err != nil {
		return nil, err
	}

	return &pb.Frame{Id: id.String()}, nil
}
