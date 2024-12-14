package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-02-bidirectional-streaming/bidirectional"
)

// creates a new 'Message'
func makeMessage(message string) *pb.Message {
	return &pb.Message{
		Message: message,
	}
}

// sends a stream of messages
func generateMessages() []*pb.Message {
	return []*pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}
}

func sendMessage(client pb.BidirectionalClient) {

	// creates a bidirectional stream
	stream, err := client.GetServerResponse(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// receives responses using goroutine
	go func() {
		for {
			resp, err := stream.Recv()

			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatalf("Error receiving message: %v", err)
			}

			time.Sleep(10 * time.Millisecond)

			log.Printf("[server to client] %s", resp.GetMessage())
		}
	}()

	// sends messages to the server (while receiving responses asynchronously)
	messages := generateMessages()

	for _, msg := range messages {
		log.Printf("[client to server] %s", msg.GetMessage())

		if err := stream.Send(msg); err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	stream.CloseSend()
}

func main() {

	// removes the default log timestamp
	log.SetFlags(0)

	// connects to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBidirectionalClient(conn)

	sendMessage(client)
}
