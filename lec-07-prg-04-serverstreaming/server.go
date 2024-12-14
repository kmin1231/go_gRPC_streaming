package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-04-serverstreaming/serverstreaming"
)

type server struct {
	pb.UnimplementedServerStreamingServer
}

func (s *server) GetServerResponse(req *pb.Number, stream pb.ServerStreaming_GetServerResponseServer) error {
	messages := []string{
		"message #1",
		"message #2",
		"message #3",
		"message #4",
		"message #5",
	}

	fmt.Printf("Server processing gRPC server-streaming {%d}.\n", req.GetValue())

	// channel to track whether all messages have been sent
	doneChannel := make(chan struct{})

	// goroutine to send messages
	go func() {
		for _, msg := range messages {
			if err := stream.Send(&pb.Message{Message: msg}); err != nil {
				log.Printf("Failed to send message: %v", err)
				close(doneChannel)
				return
			}
		}

		// all messages are sent -> closes the channel
		close(doneChannel)
	}()

	// waits for the goroutine to finish sending all messages
	<-doneChannel

	return nil
}

func main() {

	log.SetFlags(0)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterServerStreamingServer(grpcServer, &server{})

	log.Println("Starting server. Listening on port 50051.")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
