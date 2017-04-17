package game

import (
	"strings"
	"sort"
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
	MainColonist Colonist
	MainBase Base
}

type CommandParser struct{
	In Input
}

type Colonist struct {
	Name string
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
	exitCommand := func(args []string) (err error) {
		err = &ExitCalled{"A subroutine called exit"}
		return
	}
	g.CommandsMap["exit"] = exitCommand
	g.CommandsMap["x"] = exitCommand
	g.CommandsMap ["commands"] = g.PrintCommands
	g.CommandsMap ["look"] = g.Look
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
	// g.CreateBase()
	// g.CreateColonist()
	err = nil
	return
}

func (g *GameManager) CreateColonist() {
	out := g.Out
	out.Println("What would you like to name your first colonist?")

	in := g.In
	name := in.Read()
	c := g.CreateColonistWithName(name)
	g.MainColonist = c
	out.Println("Creating a colonist with the name: \"" + c.Name + "\", is this correct? (y/n)")
	answer := strings.ToLower(in.Read())
	if (answer[0] == 'y') {
		out.Println("Colonist with name \"" + c.Name + "\" created.")
	} else {
		g.CreateColonist()
	}
	return
}

func (g *GameManager) CreateColonistWithName(name string) (c Colonist){
	c = Colonist{Name : name}
	return
}

func (g *GameManager) CreateBase() {
	out := g.Out
	out.Println("What would you like to name your base?")

	in := g.In
	name := in.Read()
	b := g.CreateBasesWithName(name)
	g.MainBase = b
	out.Println("Naming your base: \"" + b.Name + "\", is this correct? (y/n)")
	answer := strings.ToLower(in.Read())
	if (answer[0] == 'y') {
		out.Println("Base with name \"" + b.Name + "\" created.")
	} else {
		g.CreateBase()
	}
	return
}

type Base struct {
	Name string
}

func (g *GameManager) CreateBasesWithName(name string) (b Base){
	b = Base{Name : name}
	return
}

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
	SHOW_POSSIBILITIES
	CANCEL
	RETRY
)

func (g *GameManager) Look(args []string) (err error) {
	out := g.Out
	in := g.In
	var dir Direction

	if (len(args) <= 1) {
		for dir = RETRY ; (dir == RETRY); {
			out.Println("In which direction would you like to look? ('x' to cancel or 'd' for possible directions)")
			dirString := strings.ToLower(in.Read())
			if (strings.HasPrefix(dirString, "look")) {
				dirString = strings.TrimSpace(dirString[4:])
			}
			dir = g.GetDirection(dirString)
			if (dir == RETRY) {
				out.Println("I do not know what direction " + dirString + " is.")
				out.Println("The possible directions are N, E, S, W")
			}
		}
	} else {
		dirString := args[1]
		dir = g.GetDirection(dirString)
	}

	switch dir {
	case NORTH:
		out.Println("You looked North!")
	case EAST:
		out.Println("You looked East!")
	case SOUTH:
		out.Println("You looked South!")
	case WEST:
		out.Println("You looked West!")
	case SHOW_POSSIBILITIES:
		out.Println("The possible directions are N, E, S, W")
	case CANCEL:
		out.Println("You cancelled looking")
	default:
		out.Println("I don't know which way you looked")
	}

	err = nil
	return
}

func  (g *GameManager) GetDirection(dirString string) (dir Direction) {

	dirChar := dirString[0]
	switch dirChar {
	case 'n':
		dir = NORTH
	case 'e':
		dir = EAST
	case 's':
		dir = SOUTH
	case 'w':
		dir = WEST
	case 'x':
		dir = CANCEL
	case 'd':
		dir = SHOW_POSSIBILITIES
	default:
		dir = RETRY
	}
	return
}