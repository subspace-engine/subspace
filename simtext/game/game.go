package game

import (
	"strings"
	"sort"
	"strconv"
)

type Output interface  {
	Print(s string)
	Println(s string)
}

type Input interface  {
	Read() (s string)
}

type GameManager struct{
	Out Output
	In Input
	CommandsMap map[string]func(args []string) error
	World *World
}

type CommandParser struct{
	In Input
}

func (g *GameManager) Start() {
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

func (g *GameManager) InitializeCommandsMap() {
	g.CommandsMap = make(map[string]func(args []string) error)
	g.CommandsMap["exit"] = g.Exit
	g.CommandsMap["x"] = g.Exit
	g.CommandsMap["commands"] = g.PrintCommands
	g.CommandsMap["look"] = g.Look
	g.CommandsMap["terrain"] = g.DrawTerrain
}

func (g *GameManager) MainLoop() {
	g.StartNewGame()
	Loop:
	for {
		if doExit := g.LoopStep() ; doExit != nil {
			break Loop
		}
	}
	g.Out.Println("Exiting game. Goodbye!")
}

func (g *GameManager) DrawTerrain(args []string) (err error) {
	z := 2
	if (len(args) > 1) {
		z, err = strconv.Atoi(strings.TrimSpace(args[1]))
	}
	out := g.Out

	out.Println("Drawing terrain")
	drawnTerrain, _ := g.World.Terrain.DrawnTerrainAtZ(z)
	out.Print(drawnTerrain)
	return nil
}


func (g *GameManager) Exit(args []string) (err error) {
	in := g.In
	out := g.Out
	out.Println("Are you sure you want to exit? (y/n)")
	answer := strings.ToLower(in.Read())
	if (answer[0] == 'y') {
		out.Println("Thanks for playing!")
		err = &ExitCalled{"A subroutine called exit"}
	} else {
		out.Println("Cancelled exit")
		err = nil
	}
	return
}


func (g *GameManager) LoopStep() (err error) {
	in := g.In
	out := g.Out
	err = nil

	line := strings.Fields(strings.ToLower(in.Read()))
	command := line[0]
	commandsMap := g.CommandsMap

	commandFunction := commandsMap[command]
	if commandFunction == nil {
		out.Println("Command \"" + command + "\" not recognized.")
		g.PrintCommands([]string{})
	} else {
		err = commandFunction(line)
	}
	return
}

func (g *GameManager) PrintCommands(args []string) (err error) {
	commandsMap := g.CommandsMap

	keys := make([]string, 0, len(commandsMap))
	for k := range commandsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	g.Out.Print("The available commands are: ")
	g.Out.Println(strings.Join(keys, ", "))

	err = nil
	return
}

func (g *GameManager) StartNewGame() (err error) {
	out := g.Out
	out.Println("Starting a new game")
	g.CreateWorld()
	err = nil
	return
}

