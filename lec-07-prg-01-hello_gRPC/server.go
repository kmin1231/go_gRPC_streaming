package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	// imports gRPC module
	"google.golang.org/grpc"

	// imports *.proto file
	// [Go] NO need to import auto-generated proto source code separately
	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-01-hello_gRPC/hello_gRPC"

	hello_gRPC "github.com/kmin1231/simple_grpc_go/lec-07-prg-01-hello_gRPC/hello_gRPC"
)

// defines the gRPC server
type server struct {
	pb.UnimplementedMyServiceServer // required in Go (inherits the base gRPC class)
}

func (s *server) MyFunction(ctx context.Context, in *pb.MyNumber) (*pb.MyNumber, error) {
	var result int

	// 'future.ThreadPoolExecutor' in Python -> 'goroutine' in Go
	go func() {
		result = hello_gRPC.MyFunc(int(in.Value))
	}()

	// waits for input
	time.Sleep(200 * time.Millisecond)

	return &pb.MyNumber{Value: int32(result)}, nil
}

func main() {
	// listens for the gRPC server on a TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// creates a gRPC server
	s := grpc.NewServer()

	// registers the service to the gRPC server
	pb.RegisterMyServiceServer(s, &server{})

	// starts the server & listens on the specified port 50051
	fmt.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
