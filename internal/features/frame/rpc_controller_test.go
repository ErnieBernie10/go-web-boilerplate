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
