package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-03-clientstreaming/clientstreaming"
)

func makeMessage(message string) *pb.Message {
	return &pb.Message{
		Message: message,
	}
}

func generateMessages() []*pb.Message {
	return []*pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}
}

func sendMessages(client pb.ClientStreamingClient) {

	stream, err := client.GetServerResponse(context.Background())
	if err != nil {
		log.Fatalf("could not start stream: %v", err)
	}

	messages := generateMessages()

	for _, msg := range messages {
		fmt.Printf("[client to server] %s\n", msg.Message)

		if err := stream.Send(makeMessage(msg.Message)); err != nil {
			log.Fatalf("could not send message: %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}

	fmt.Printf("[server to client] %d\n", response.GetValue())
}

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewClientStreamingClient(conn)

	sendMessages(client)
}
