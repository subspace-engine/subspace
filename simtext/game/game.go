package game

import "strings"

type Output interface  {
	Print(s string)
	Println(s string)
}

type Input interface  {
	Read() (s string)
}

type LoopHandler struct{
	Out Output
	In Input
}

func (g *LoopHandler) Start() {
	g.MainLoop()
}

/*
func (g *LoopHandler) MainLoop() {
	in := g.In
	out := g.Out

	inChan := make(chan string)
	outChan := make(chan string)

	Loop:
	for {
		go func() {
			inChan <- in.Read()
		}()

		go func() {
			outChan <- "Welcome to Simtext!"
		}()

		select {
			case inMessage := <-inChan:
				exitMessage := "exit"
				if (inMessage == exitMessage) {
					out.Println("Exiting")
					break Loop
				} else {
					out.Println("received: " + inMessage + "!")
				}

			case outMessage := <-outChan:
				out.Println(outMessage)
		}
	}
}
*/

func (g *LoopHandler) MainLoop() {
	out := g.Out

	out.Println("Welcome to Simtext!")

	Loop:
	for {
		if doExit := g.LoopStep() ; doExit {
			break Loop
		}
	}
}

func (g *LoopHandler) LoopStep() (doExit bool) {
	doExit = false

	in := g.In

	command := strings.ToLower(in.Read())

	if (command == "exit") {
		doExit = true
	}
	return
}