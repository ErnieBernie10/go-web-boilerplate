package database

import (
	"context"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbContainer *postgres.PostgresContainer

var userID uuid.UUID

func mustStartPostgresContainer() (func(context.Context) error, error) {
	var (
		dbName = "framer"
		dbPwd  = "Admin123"
		dbUser = "postgres"
	)
	var err error

	dbContainer, err = postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	_, err = dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	_, err = dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	connStr := dbContainer.MustConnectionString(context.Background()) + "sslmode=disable"
	log.Println(connStr)
	u, _ := url.Parse(connStr)
	db := dbmate.New(u)

	db.MigrationsDir = []string{"../../db/migrations"}
	log.Println(db.MigrationsDir)

	err = db.CreateAndMigrate()
	if err != nil {
		return dbContainer.Terminate, err
	}

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	NewDb(dbContainer.MustConnectionString(context.Background()))
	Seed()
	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	err := NewDb(dbContainer.MustConnectionString(context.Background()))
	if err != nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {

	stats := Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestGetFrames(t *testing.T) {
	ctx := context.Background()

	_, err := Service.SaveFrame(ctx, SaveFrameParams{
		Title:       "Test",
		Description: "Test",
		UserID:      userID,
		FrameStatus: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := Service.GetFrames(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatal("Res is not length 1")
	}
}
