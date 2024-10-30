package testing

import (
	"context"
	"framer/internal/pkg"
	"framer/internal/pkg/database"
	"framer/internal/pkg/view"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc/metadata"
)

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

	db.AutoDumpSchema = false
	err = db.CreateAndMigrate()
	if err != nil {
		return dbContainer, err
	}

	return dbContainer, err
}

func SetupTests(m *testing.M) (*database.DbService, uuid.UUID, func(context.Context) error) {
	dbContainer, err := MustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	db, err := database.NewDb(dbContainer.MustConnectionString(context.Background()))
	if err != nil {
		log.Fatalln("could create db", err)
	}

	userID, err := database.Seed(db)
	if err != nil {
		log.Fatalf("could not seed user: %v", err)
	}

	return db, userID, dbContainer.Terminate
}

func WithMockUser(ctx context.Context, userId uuid.UUID) (context.Context, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &pkg.Claims{
		Email: "john@doe.com",
		ID:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	refreshClaim := &pkg.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)

	tokenString, err := token.SignedString(pkg.JwtSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refresh.SignedString(pkg.JwtSecret)
	if err != nil {
		return nil, err
	}

	ctx = view.ContextWithToken(ctx, func() (string, string) {
		return tokenString, refreshTokenString
	})

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	return ctx, nil
}
