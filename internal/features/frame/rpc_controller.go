package frame

import (
	"context"
	"database/sql"
	"errors"
	"framer/internal/pkg"
	"framer/internal/pkg/database"
	"framer/internal/pkg/rpc"
	pb "framer/internal/proto"

	"github.com/google/uuid"
)

type FrameController struct {
	pb.FrameServiceServer
	Db *database.DbService
}

// DeleteFrame implements proto.FrameServiceServer.
func (c *FrameController) DeleteFrame(ctx context.Context, req *pb.DeleteByIdRequest) (*pb.EmptyResponse, error) {
	claims, err := rpc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = c.Db.Queries.DeleteFrame(ctx, database.DeleteFrameParams{
		ID:     id,
		UserID: claims.ID,
	})

	if err != nil {
		return nil, err
	}

	return &pb.EmptyResponse{}, nil
}

// ListFrames implements proto.FrameServiceServer.
func (c *FrameController) ListFrames(ctx context.Context, req *pb.EmptyResponse) (*pb.ListFramesResponse, error) {
	claims, err := rpc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	frames, err := c.Db.Queries.GetFrames(ctx, claims.ID)
	if err != nil {
		return nil, err
	}

	res := &pb.ListFramesResponse{
		Frames: make([]*pb.Frame, len(frames)),
	}

	for i, frame := range frames {
		res.Frames[i] = &pb.Frame{
			Id:          frame.ID.String(),
			Title:       frame.Title,
			Description: frame.Description,
			Content:     frame.Content,
			FrameStatus: pb.FrameStatus(frame.FrameStatus),
			ContentType: pb.ContentType(frame.ContentType),
		}

		if frame.FileID.Valid {
			res.Frames[i].File.Id = frame.FileID.UUID.String()
			res.Frames[i].File.FileName = frame.FileName.String
		}
	}
	return res, nil
}

// UpdateFrame implements proto.FrameServiceServer.
func (c *FrameController) UpdateFrame(ctx context.Context, req *pb.UpdateFrameRequest) (*pb.Frame, error) {
	claims, err := rpc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	frame, err := create(req.Title, req.Description, claims.ID, uuid.NullUUID{})
	if err != nil {
		return nil, err
	}

	uow, err := c.Db.NewUnitOfWork()
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	id, err = SaveFrameWithFile(ctx, uow, frame, req.File.Content, req.File.FileName)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	uow.Commit()

	return &pb.Frame{Id: id.String()}, nil
}

func (c *FrameController) CreateFrame(ctx context.Context, req *pb.CreateFrameRequest) (*pb.Frame, error) {
	claims, err := rpc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	frame, err := create(req.Title, req.Description, claims.ID, uuid.NullUUID{})
	if err != nil {
		return nil, err
	}

	uow, err := c.Db.NewUnitOfWork()
	defer uow.Rollback()
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	if req.File != nil {
		id, err = SaveFrameWithFile(ctx, uow, frame, req.File.Content, req.File.FileName)
	} else {
		id, err = SaveFrame(ctx, uow, frame)
	}

	if err != nil {
		return nil, err
	}

	uow.Commit()

	return &pb.Frame{Id: id.String()}, nil
}

func (c *FrameController) GetFrame(ctx context.Context, req *pb.GetByIdRequest) (*pb.Frame, error) {
	claims, err := rpc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	frame, err := c.Db.Queries.GetFrame(ctx, database.GetFrameParams{
		ID:     id,
		UserID: claims.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(pkg.ErrNotFound, errors.New("frame not found"))
		}
		return nil, err
	}

	res := &pb.Frame{
		Id:          frame.ID.String(),
		Title:       frame.Title,
		Description: frame.Description,
		Content:     "",
		FrameStatus: pb.FrameStatus(frame.FrameStatus),
		ContentType: pb.ContentType(frame.ContentType),
	}

	if frame.FileID.Valid {
		res.File.Id = frame.FileID.UUID.String()
		res.File.FileName = frame.FileName.String
	}

	return res, nil
}
