package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"postman/personpb"

	"google.golang.org/grpc"
)

// PersonServer implements the gRPC service
type PersonServer struct {
	personpb.UnimplementedPersonServiceServer
}

// GetPerson method implementation
func (s *PersonServer) GetPerson(ctx context.Context, req *personpb.First) (*personpb.First, error) {
	log.Printf("Received request for: %s", req.Name)

	// Returning the same person details for simplicity
	return &personpb.First{
		Name: req.Name,
		Details:  req.Details + 5,
	}, nil
}

func main() {
	// Create a TCP listener
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	personpb.RegisterPersonServiceServer(grpcServer, &PersonServer{})

	fmt.Println("gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
