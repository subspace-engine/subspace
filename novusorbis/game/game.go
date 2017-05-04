package game

import (
	"strings"
	"sort"
	"strconv"
	"time"
)

const GAME_LOGO = `
             :::;;;;;;::,
        ,,,:::;;;;;:""'''":,
      .,,,:::;;;;"          '
    ...,,,:::;;;"  N O V U S "
   ...,,,,:::;;;"             "
    ...,,,:::;;;". O R B I S  ";
    ...,,,::::;;;".         .";;;
    ....,,,::::;;;;;:,....,:;;;;:
     ....,,,:::::;;;;;;;;;;;;;:::
      ....,,,,::::::;;;;;;;::::::
        ....,,,,,:::::::::::::,,,
         .....,,,,,,,,,,,,,,,,,,
            .......,,,,,,,,,...
               ...............




`

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
	LetterToDirection map[rune]Direction
	DirectionToString map[Direction]string
}

type CommandParser struct{
	In Input
}

func (g *GameManager) Start() {
	g.PrintLogo()
	g.InitializeCommandsMap()
	g.SetUpDirectionMaps()
	g.MainLoop()
}

func (g *GameManager) PrintLogo() {
	out := g.Out
	for _, r := range GAME_LOGO {
		c := string(r)
		out.Print(c)
		if (r == rune('\n')) {
			time.Sleep(250 * time.Millisecond)
		}
	}
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
	g.CommandsMap["q"] = g.Exit
	g.CommandsMap["commands"] = g.PrintCommands
	g.CommandsMap["look"] = g.Look
	g.CommandsMap["move"] = g.Move
	g.CommandsMap["pos"] = g.Position
	g.CommandsMap["draw"] = g.DrawWorld
}

func (g *GameManager) MainLoop() {
	g.StartNewGame()
	Loop:
	for {
		if doExit := g.LoopStep() ; doExit != nil {
			break Loop
		}
	}
}

func (g *GameManager) DrawWorld(args []string) (err error) {
	z := g.World.Size/2
	if (len(args) > 1) {
		z, err = strconv.Atoi(strings.TrimSpace(args[1]))
	}
	out := g.Out

	out.Println("Drawing terrain at " + strconv.Itoa(z))
	drawnTerrain, _ := g.World.DrawnWorldAtZ(z)

	out.Println(drawnTerrain)

	drawnTerrain, _ = g.World.DrawnWorldAtZ(z-1)

	out.Print(drawnTerrain)
	return nil
}

func (g *GameManager) Position(args []string) (err error) {
	err = nil
	p := g.World.MainColonist.Avatar.Position()
	g.Out.Println("(" + strconv.Itoa(p.x) + ", " + strconv.Itoa(p.y) + ", " + strconv.Itoa(p.z) + ")")
	return
}

func (g *GameManager) Exit(args []string) (err error) {
	in := g.In
	out := g.Out
	out.Println("Are you sure you want to exit? (y/n)")
	answer := strings.ToLower(in.Read())
	if (len(answer) > 0 && answer[0] == 'y') {
		out.Println("Returning to reality.")
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

	if len(line) == 0 {
		return
	}

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


