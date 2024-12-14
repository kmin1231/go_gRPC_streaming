package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-04-serverstreaming/serverstreaming"
)

func receiveMessages(stub pb.ServerStreamingClient) {

	req := &pb.Number{Value: 5}

	stream, err := stub.GetServerResponse(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling GetServerResponse: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {

			// prevents EOF error
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error receiving message: %v", err)
		}
		fmt.Printf("[server to client] %s\n", resp.Message)
	}
}

func main() {

	log.SetFlags(0)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewServerStreamingClient(conn)

	receiveMessages(client)
}
