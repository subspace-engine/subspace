package game

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

		for i := 0; i < 1; i++ {
			select {
			case inMessage := <-inChan:
				exitMessage := "exit"
				if (inMessage == exitMessage) {
					out.Println("Exiting")
					break Loop
				} else {
					out.Println("received: " + inMessage + "!")
				}

			// case outMessage := <-outChan:
				// out.Println("outMessage: " + outMessage)
			}
		}
	}
}