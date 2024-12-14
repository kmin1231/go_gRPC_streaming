package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-03-clientstreaming/clientstreaming"
)

type server struct {
	pb.UnimplementedClientStreamingServer
}

func (s *server) GetServerResponse(stream pb.ClientStreaming_GetServerResponseServer) error {

	log.Println("Server processing gRPC client-streaming.")

	// counts the number of received messages
	count := 0

	// creates a channel to handle counting and sending responses
	done := make(chan bool)

	// goroutine to receive messages asynchronously
	go func() {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				done <- false
				return
			}

			// increments count for each message received
			count++
		}
	}()

	// waits for the goroutine to finish
	if <-done {
		return stream.SendAndClose(&pb.Number{Value: int32(count)})
	} else {
		return fmt.Errorf("error receiving messages")
	}
}

func main() {

	log.SetFlags(0)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterClientStreamingServer(s, &server{})

	log.Println("Starting server. Listening on port 50051.")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
