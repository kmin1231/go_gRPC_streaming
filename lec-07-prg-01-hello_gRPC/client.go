package main

import (
	"context"
	"fmt"
	"log"

	// imports gRPC module
	"google.golang.org/grpc"

	// imports *.proto file
	// [Go] NO need to import auto-generated proto source code separately
	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-01-hello_gRPC/hello_gRPC"
)

func main() {
	// creates an insecure gRPC communication channel
	conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()

	// generates a gRPC client stub
	client := pb.NewMyServiceClient(conn)

	// creates a request message (for remote function call)
	request := &pb.MyNumber{Value: 4}

	// calls the remote function 'MyFunction'
	response, err := client.MyFunction(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to call MyFunction: %v", err)
	}

	// prints the result
	fmt.Printf("gRPC result: %d\n", response.GetValue())
}
