package main

import (
	"context"
	"framer/internal/features/auth"
	"framer/internal/features/frame"
	"framer/internal/pkg/database"
	"framer/internal/pkg/logger"
	pb "framer/internal/proto"
	"log"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	port := os.Getenv("GRPC_PORT")

	server := NewServer(ctx)

	defer Shutdown()

	lis, err := net.Listen("tcp", ":"+port)
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("grpc server listening on port %s", lis.Addr().String())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewServer(ctx context.Context) *grpc.Server {
	zapLogger := logger.NewLogger()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zapLogger),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zapLogger),
		)),
	)

	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := database.NewDb(databaseUrl)
	if err != nil {
		panic(err.Error())
	}

	pb.RegisterFrameServiceServer(grpcServer, &frame.FrameController{Db: db})
	pb.RegisterAppUserServiceServer(grpcServer, &auth.AuthController{Db: db})

	return grpcServer
}

func Shutdown() {
	database.Shutdown()
}
