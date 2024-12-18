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
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *DbService

var userID uuid.UUID

func MustStartPostgresContainer() (*postgres.PostgresContainer, error) {
	var dbContainer *postgres.PostgresContainer
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
		return dbContainer, err
	}

	_, err = dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer, err
	}

	connStr := dbContainer.MustConnectionString(context.Background()) + "sslmode=disable"
	log.Println(connStr)
	u, _ := url.Parse(connStr)
	db := dbmate.New(u)

	db.MigrationsDir = []string{"../../../db/migrations"}
	log.Println(db.MigrationsDir)

	err = db.CreateAndMigrate()
	if err != nil {
		return dbContainer, err
	}

	return dbContainer, err
}

func TestMain(m *testing.M) {
	dbContainer, err := MustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	db, err = NewDb(dbContainer.MustConnectionString(context.Background()))
	if err != nil {
		log.Fatalln("could create db", err)
	}

	userID, err = Seed(db)
	if err != nil {
		log.Fatalf("could not seed user: %v", err)
	}
	m.Run()

	if dbContainer.Terminate(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestHealth(t *testing.T) {

	stats := db.Health()

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

	_, err := db.Queries.SaveFrame(ctx, SaveFrameParams{
		Title:       "Test",
		Description: "Test",
		UserID:      userID,
		FrameStatus: 1,
		ID:          uuid.New(),
		ContentType: int16(1),
		Content:     "",
		FileID:      uuid.NullUUID{},
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := db.Queries.GetFrames(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatal("Res is not length 1")
	}
}
