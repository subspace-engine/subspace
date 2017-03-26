package main

import "fmt"
import "bufio"
import "net"
import "strings"
import "subspace/world"

func handleConnection(con net.Conn) {
	defer con.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
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
		fmt.Fprintf(rw, "Read : %s of length %d\n", s, len(s))
		rw.Flush()
	}
}

func main() {
	obj := world.NewObject()
	fmt.Println(obj)
	fmt.Println("Subspace Core Starting")
	ln, err := net.Listen("tcp", ":4444")
	if err != nil {
		fmt.Println("Unable to create listener!")
		return
	}
	for {
		con, err := ln.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection!")
			return
		}
		go handleConnection(con)
	}
}
