package main

import (
	"fmt"
	"log"

	pb "github.com/rynkruger/subspace/clientserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	fmt.Println("Starting client!")
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // TODO Make secure

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGameClient(conn)

	stream, err := client.Play(context.Background())

	if err != nil {
		log.Fatalf("Could not start client stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive info : %v", err)
			}

			fmt.Println("Got info: %v", in)
		}
	}()

	commands := []*pb.Command{
		{Text:"firstCommamnd"},
		{Text:"secondCommand"},
	}

	for _, command := range commands {
		if err := stream.Send(command); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}

	err = stream.CloseSend()

	if err != nil {
		log.Fatalf("%v.CloseSend() got error %v, want %v", stream, err, nil)
	}
	<-waitc
}

