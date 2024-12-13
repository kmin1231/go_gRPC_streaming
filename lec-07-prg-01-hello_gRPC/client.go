package main

import (
	"context"
	"fmt"
	"log"

	// (1) imports gRPC module
	"google.golang.org/grpc"

	// (2) imports *.proto file
	// in Golang, no need to import auto-generated proto source code
	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-01-hello_gRPC/hello_gRPC"
)

func main() {
	// (3) creates an insecure gRPC communication channel
	conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()

	// (4) generates a gRPC client stub
	client := pb.NewMyServiceClient(conn)

	// (5) creates a request message (for remote function call)
	request := &pb.MyNumber{Value: 4}

	// (6) calls the remote function
	response, err := client.MyFunction(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to call MyFunction: %v", err)
	}

	// (7) prints the result
	fmt.Printf("gRPC result: %d\n", response.GetValue())
}
