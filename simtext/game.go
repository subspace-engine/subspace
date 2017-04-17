package simtext

type Output interface  {
	Print(s string)
	Println(s string)
}

type Input interface  {
	Read() (s string)
}

type Game struct{
	Out Output
	In Input
}

func (g *Game) Start() {
	g.MainLoop()
}

func (g *Game) MainLoop() {
	in := g.In
	out := g.Out

	inChan := make(chan string)
	outChan := make(chan string)

	for true {
		go func() {
			inChan <- in.Read()
		}()

		go func() {
			outChan <- "Welcome to Simtext!"
		}()

		for i := 0; i < 2; i++ {
			select {
			case msg1 := <-inChan:
				out.Println("received " + msg1)
			case msg2 := <-outChan:
				out.Println("received " + msg2)
			}
		}
	}
}