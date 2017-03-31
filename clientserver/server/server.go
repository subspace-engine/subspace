package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/rynkruger/subspace/clientserver"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement GreeterServer.
type server struct{}

func (s *server) SendMessages(in pb.Messenger_SendMessagesServer)  error {
	command, err := in.Recv()
	if err != nil {
		log.Fatalf("Failed to receive anything: %v", err)
	}
	fmt.Println("Received command: %v", command)

	command, err = in.Recv()
	if err != nil {
		log.Fatalf("Failed to receive anything: %v", err)
	}
	fmt.Println("Received command: %v", command)
	fmt.Println("Done receiving messages.")

	in.SendAndClose(&pb.Response{"hello"})
	return nil
}

func main() {
	fmt.Println("Starting server!")
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessengerServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	fmt.Println("Closing server")
}
