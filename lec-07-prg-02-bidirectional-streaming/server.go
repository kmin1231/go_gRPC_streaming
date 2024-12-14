package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/kmin1231/simple_grpc_go/lec-07-prg-02-bidirectional-streaming/bidirectional"
)

type server struct {
	pb.UnimplementedBidirectionalServer
}

func (s *server) GetServerResponse(stream pb.Bidirectional_GetServerResponseServer) error {

	log.Println("Server processing gRPC bidirectional streaming.")

	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}

			go func(req *pb.Message) {
				if stream.Context().Err() != nil {
					// If the context is canceled, log and return
					log.Printf("Context canceled while sending: %v", stream.Context().Err())
					return
				}

				if err := stream.Send(req); err != nil {
					log.Printf("Error sending message: %v", err)
				}
			}(req)
		}
	}()

	// keeps the server running
	select {}
}

func main() {

	// removes timestamp for log messages
	log.SetFlags(0)

	// setting for gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterBidirectionalServer(s, &server{})

	// starts listening for incoming gRPC requests on port 50051
	fmt.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
