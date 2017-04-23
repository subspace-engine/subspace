package server

import (
	"fmt"
	engine "github.com/subspace-engine/subspace/engine"
	"github.com/subspace-engine/subspace/world"
	"google.golang.org/grpc"
	"io"
	"net"
)

type server struct {
	done     chan struct{}
	game     world.PlayerWorld
	listener net.Listener
	grpc     *grpc.Server
}

func NewServer(port int, done chan struct{}) *server {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}
	w := world.NewTestWorld()
	go w.Run(done)
	return &server{done, w, ln, grpc.NewServer()}
}

func (s *server) Run() {
	engine.RegisterMessageStreamServer(s.grpc, s)
	s.grpc.Serve(s.listener)
}

func (s *server) Close() {
	close(s.done)
	s.grpc.GracefulStop()
	s.listener.Close()
}

func (s *server) UpdateStream(stream engine.MessageStream_UpdateStreamServer) error {
	player := s.game.RequestPlayer()
	updater := s.game.Updater(player.Id())
	recvErr := make(chan error)
	go func() {
		actor := s.game.Actor()
		for {
			act, err := stream.Recv()
			if err == io.EOF {
				fmt.Printf("Connection closed by client %d", player.Id())
			}
			if err != nil {
				fmt.Printf("Client %d received error %s", player.Id(), err.Error())
				recvErr <- err
			}
			fmt.Printf("server action : %v", act)
			actor <- engine.PlayerAction{player.Id(), act}
		}
	}()
	stream.Send(&engine.ServerUpdate{"Fear the boom! Only some have the ability to nutralise it."})
	var err error
	for {
		select {
		case <-s.done:
			stream.Send(&engine.ServerUpdate{fmt.Sprintf("Client %d, server is closing", player.Id())})
			return nil
		case err = <-recvErr:
			return err
		case up := <-updater:
			stream.Send(up)
		}
	}
	return nil
}
