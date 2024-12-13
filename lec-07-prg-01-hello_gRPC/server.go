package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// (1) imports gRPC module
	"google.golang.org/grpc"

	// (2) imports *.proto file
	// [Golang] NO need to import auto-generated proto source code separately
	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-01-hello_gRPC/hello_gRPC"

	hello_gRPC "github.com/kmin1231/simple_grpc_go/lec-07-prg-01-hello_gRPC/hello_gRPC"
)

// (3) defines the gRPC server
type server struct {
	pb.UnimplementedMyServiceServer // required in Go (inherit the base gRPC class)
}

// // (4) implements the RPC function defined in proto file
func (s *server) MyFunction(ctx context.Context, in *pb.MyNumber) (*pb.MyNumber, error) {
	// uses the custom function from hello_grpc.go
	result := hello_gRPC.MyFunc(int(in.Value))
	return &pb.MyNumber{Value: int32(result)}, nil
}

func main() {
	// (6) listens for the gRPC server on a TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// (7) creates a gRPC server
	s := grpc.NewServer()

	// (8) registers the service to the gRPC server
	pb.RegisterMyServiceServer(s, &server{})

	// (9) starts the server & listens on the specified port 50051
	fmt.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
