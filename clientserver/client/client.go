package main

import (
	"fmt"
	"log"

	pb "github.com/rynkruger/subspace/clientserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	fmt.Println("Starting client!")
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // TODO make secure?

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewMessengerClient(conn)

	stream, err := client.SendMessages(context.Background())
	if err != nil {
		log.Fatalf("Could not start client stream: %v", err)
	}

	stream.Send(&pb.Command{"command1"})
	stream.Send(&pb.Command{"command2"})
	response, err:= stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Could not receive response from server: %v", err)
	}
	log.Printf("Server response: %s", response)
}

