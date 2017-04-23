package main

import (
	"github.com/subspace-engine/subspace/world"
	"bufio"
	"net"
	"strings"
)

func handleConnection(con net.Conn, sim *world.Sim) {
	defer con.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	agent := world.NewWObject("you", world.NewBox(0,0,0,1,2,1))
	client := world.NewClient(agent)
	sim.AddClient(client)
	rw.WriteString("Subspace 1.0\n")
	rw.Flush()
	for {
		s, err := rw.ReadString('\n')
		if err != nil {
			return
		}
		s = strings.TrimSpace(s)
		if s == "done" {
			rw.WriteString("bye\n")
			rw.Flush()
			return
		}
		client.Queue <- s
		resp := <- client.Queue
				rw.WriteString(resp)
		rw.Flush()
	}
}

func RunSim(sim *world.Sim) {
	world1 := world.NewWObject("world", world.NewBox(-1000000, -1000000, -1000000, 2000000, 2000000, 2000000))
	sim.AddWorld(world1)
	sim.Run()
}

func main() {
	runTermbox()

	/*
	obj := world.NewObject()
	fmt.Println(obj)
	fmt.Println("Subspace Core Starting")
	ln, err := net.Listen("tcp", ":4444")
	if err != nil {
		fmt.Println("Unable to create listener!")
		return
	}
	sim := world.NewSim()
	go RunSim(&sim)
	for {
		con, err := ln.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection!")
			return
		}

		go handleConnection(con, &sim)
	}

	*/
}
