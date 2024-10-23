package main

import (
	"framer/internal/database"
	"framer/internal/features/frame"
	pb "framer/internal/proto"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterFrameServiceServer(grpcServer, &frame.FrameController{})

	databaseUrl := os.Getenv("DATABASE_URL")
	database.NewDb(databaseUrl)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("grpc server listening on port %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
