package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/subspace-engine/subspace/engine"
	"google.golang.org/grpc"
	"io"
	"os"
	"sync"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:11111", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	stream, err := engine.NewMessageStreamClient(conn).UpdateStream(context.Background())
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			update, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Connection closed by server")
			}
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", update.Desc)
		}
	}()

	go func() {
		defer wg.Done()

	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("client: %s\n", text)
		stream.Send(&engine.Action{text})
		if text == "exit" {
			time.Sleep(time.Second)
			break
		}
	}
	fmt.Println("Exiting now")
	wg.Wait()
	conn.Close()
}
