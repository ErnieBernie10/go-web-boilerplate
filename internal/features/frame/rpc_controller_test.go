package frame

import (
	"context"
	"errors"
	"framer/internal/pkg"
	"framer/internal/pkg/database"
	test "framer/internal/pkg/testing"
	"framer/internal/proto"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/google/uuid"
)

var db *database.DbService
var userId uuid.UUID
var teardown func(context.Context) error

func TestMain(m *testing.M) {
	db, userId, teardown = test.SetupTests(m)
	m.Run()

	if teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container")
	}
}

func TestCreateFrame_ReturnsFrame(t *testing.T) {
	c := &FrameController{
		Db: db,
	}
	ctx := context.Background()
	ctx, err := test.WithMockUser(ctx, userId)
	if err != nil {
		t.Fail()
	}

	frame, err := c.CreateFrame(ctx, &proto.CreateFrameRequest{
		Title:       "Test",
		Description: "Test",
		Content:     "Test",
		FrameStatus: proto.FrameStatus_Active,
	})

	if err != nil {
		t.Error("Error creating frame")
	}

	if frame.Id == "" {
		t.Error("Frame ID should not be empty")
	}
}

func TestGetFrame_ReturnsFrame(t *testing.T) {
	ctx := context.Background()
	c := &FrameController{
		Db: db,
	}

	id, err := db.Queries.SaveFrame(ctx, database.SaveFrameParams{
		ID:          uuid.New(),
		Title:       "Test",
		Description: "Test",
		UserID:      userId,
	})

	if err != nil {
		t.Fail()
	}

	ctx, err = test.WithMockUser(ctx, userId)
	if err != nil {
		t.Error("Could not create mock user")
	}

	frame, err := c.GetFrame(ctx, &proto.GetByIdRequest{
		Id: id.String(),
	})

	if err != nil {
		t.Error("Error getting frame")
	}

	if frame.Id == "" {
		t.Error("Frame ID should not be empty")
	}
}

func TestGetFrame_ReturnsErrorIfWrongUser(t *testing.T) {
	c := &FrameController{
		Db: db,
	}
	ctx := context.Background()

	id, err := db.Queries.SaveFrame(ctx, database.SaveFrameParams{
		ID:          uuid.New(),
		Title:       "Test",
		Description: "Test",
		UserID:      userId,
	})

	ctx, err = test.WithMockUser(ctx, uuid.New())
	if err != nil {
		t.Fail()
	}

	_, err = c.GetFrame(ctx, &proto.GetByIdRequest{
		Id: id.String(),
	})

	if err == nil {
		t.Fail()
	}

	if !errors.Is(err, pkg.ErrNotFound) {
		t.Fatal("Error should be ErrNotFound but was", err)
	}
}

func TestGetFrame_ReturnsNotFoundError(t *testing.T) {
	c := &FrameController{
		Db: db,
	}

	ctx := context.Background()
	ctx, err := test.WithMockUser(ctx, userId)
	if err != nil {
		t.Fail()
	}

	_, err = c.GetFrame(ctx, &proto.GetByIdRequest{
		Id: uuid.NewString(),
	})

	if err == nil {
		t.Fail()
	}

	if !errors.Is(err, pkg.ErrNotFound) {
		t.Fatal("Error should be ErrNotFound but was", err)
	}
}
