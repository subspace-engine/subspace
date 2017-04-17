package game

import "strings"

type Output interface  {
	Print(s string)
	Println(s string)
}

type Input interface  {
	Read() (s string)
}

type LoopRunner struct{
	Out Output
	In Input
	CommandsMap map[string]func() error
}

type CommandParser struct{
	In Input
}

type Colonist struct {
	Name string
}

type Game struct {
	MainColonist Colonist
}

func (g *LoopRunner) Start() {
	out := g.Out
	out.Println("Welcome to Simtext!")
	g.InitializeCommandsMap()
	g.MainLoop()
}

type ExitCalled struct {
	s string
}

func (e *ExitCalled) Error() string {
	return e.s
}

func (g *LoopRunner) InitializeCommandsMap() {
	g.CommandsMap = make(map[string]func() error)
	exitCommand := func() (err error) {
		err = &ExitCalled{"A subroutine called exit"}
		return
	}
	g.CommandsMap["exit"] = exitCommand
	g.CommandsMap ["start"] = g.StartNewGame
	g.CommandsMap ["commands"] = g.PrintCommands
}

func (g *LoopRunner) MainLoop() {
	Loop:
	for {
		if doExit := g.LoopStep() ; doExit != nil {
			break Loop
		}
	}
}

func (g *LoopRunner) PrintCommands() (err error) {
	commandsMap := g.CommandsMap
	for key := range commandsMap {
		g.Out.Print(key + ",")
	}
	err = nil
	return
}

func (g *LoopRunner) LoopStep() (err error) {
	in := g.In
	out := g.Out
	err = nil

	command := strings.ToLower(in.Read())
	commandsMap := g.CommandsMap

	commandFunction := commandsMap[command]
	if commandFunction == nil {
		out.Println("Command \"" + command + "\" not recognized.")
	} else {
		err = commandFunction()
	}
	return
}

func (g *LoopRunner) StartNewGame() (err error) {
	out := g.Out
	out.Println("Starting a new game")
	g.CreateColonist()
	err = nil
	return
}

func (g *LoopRunner) CreateColonist() {
	out := g.Out
	out.Println("What would you like to name your first colonist?")

	in := g.In
	name := in.Read()
	c := g.CreateColonistWithName(name)
	out.Println("Created a colonist with the name: \"" + c.Name + "\", is this correct?")
	answer := strings.ToLower(in.Read())
	if (answer[0] == 'y') {
		out.Println("Great!")
	} else {
		out.Println("Oh, no! I guess we'll have to try again!")
	}
}

func (g *LoopRunner) CreateColonistWithName(name string) (c Colonist){
	return Colonist{Name : name}
}
