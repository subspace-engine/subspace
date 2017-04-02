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

type GameServer struct{}

func (s *GameServer) Play(stream pb.Game_PlayServer) error {
	command, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive anything: %v", err)
	}
	fmt.Println("Received command: %v", command)

	err = stream.Send(&pb.Info{Text:"hello!"})

	if (err != nil) {
		log.Fatalf("Failed to send anything: %v", err)
	}

	return nil
}

func main() {
	fmt.Println("Starting server!")
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGameServer(grpcServer, &GameServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	fmt.Println("Closing server")
	grpcServer.GracefulStop()

}
