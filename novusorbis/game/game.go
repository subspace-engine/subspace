package game

import (
	"strings"
	"sort"
	"strconv"
	"time"
	"github.com/subspace-engine/subspace/novusorbis/world"
	"github.com/subspace-engine/subspace/novusorbis/ui"
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

type GameManager struct{
	ui.InputOutput
	CommandsMap map[string]func(args []string) error
	World *world.World
	LetterToDirection map[rune]Direction
	DirectionToString map[Direction]string
	BaseFactory BaseFactory
}

func (g *GameManager) Start() {
	// g.PrintLogo() TODO
	g.InitializeCommandsMap()
	g.SetUpDirectionMaps()
	g.CreateWorld()
	g.MainLoop()
}

func (g *GameManager) PrintLogo() {
	for _, r := range GAME_LOGO {
		c := string(r)
		g.Print(c)
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
	g.Println("Drawing terrain at " + strconv.Itoa(z))
	drawnTerrain, _ := g.World.DrawnWorldAtZ(z)

	g.Println(drawnTerrain)

	drawnTerrain, _ = g.World.DrawnWorldAtZ(z-1)

	g.Print(drawnTerrain)
	return nil
}

func (g *GameManager) Position(args []string) (err error) {
	err = nil
	p := g.World.MainColonist.Avatar.Position()
	g.Println("(" + strconv.Itoa(p.X) + ", " + strconv.Itoa(p.Y) + ", " + strconv.Itoa(p.Z) + ")")
	return
}

func (g *GameManager) Exit(args []string) (err error) {
	g.Println("Are you sure you want to exit? (y/n)")
	answer, err := g.Read()
	answer = strings.ToLower(answer)
	if (len(answer) > 0 && answer[0] == 'y') {
		g.Println("Returning to reality.")
		err = &ExitCalled{"A subroutine called exit"}
	} else {
		g.Println("Cancelled exit")
		err = nil
	}
	return
}

func (g *GameManager) LoopStep() (err error) {
	err = nil
	str, err := g.Read()
	line := strings.Fields(strings.ToLower(str))

	if len(line) == 0 {
		return
	}

	command := line[0]
	commandsMap := g.CommandsMap

	commandFunction := commandsMap[command]
	if commandFunction == nil {
		g.Println("Command \"" + command + "\" not recognized.")
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

	g.Print("The available commands are: ")
	g.Println(strings.Join(keys, ", "))

	err = nil
	return
}

func (g *GameManager) CreateWorld() (err error) {
	w := &world.World{}
	w.Size = 5
	g.World = w
	mid := w.Size/2

	c := world.Position{mid,mid,mid}
	w.Cursor = c

	w.GenerateTerrain()

	thingStore := &world.MapThingStore{}
	thingStore.Initialize()
	w.Things = thingStore

	w.Structures = world.NewStructureStore()

	w.MainColonist = g.CreateDefaultColonist() // TODO
	w.MainBase = g.CreateDefaultBase() // TODO

	pos := world.Position{mid,mid,mid}
	mainAvatar := g.World.MainColonist.Avatar
	mainAvatar.SetPosition(pos)
	thingStore.AddObjectAt(mainAvatar, pos)

	baseAvatar := g.World.MainBase.Avatar
	baseAvatar.SetPosition(pos)
	thingStore.AddObjectAt(baseAvatar, pos)

	err = nil
	return
}

func (g *GameManager) CreateDefaultColonist() (mainColonist *world.Colonist) {
	mainColonist = world.NewDefaultColonist()
	return
}

func (g *GameManager) CreateDefaultBase() (base *world.Base) {
	base = world.NewDefaultBase()
	return
}


func (g *GameManager) CreateColonist() (mainColonist *world.Colonist) {
	g.Println("What would you like to name your first colonist?")
	name, _ := g.Read()

	mainColonist = &world.Colonist{Name: name, Avatar:world.NewThing("You", "@", world.Position{2,2,2})}
	g.Println("Creating a colonist with the name: \"" + mainColonist.Name + "\", is this correct? (y/n)")
	answer, _ := g.Read()
	answer = strings.ToLower(answer)
	if (len(answer) > 0 && answer[0] == 'y') {
		g.Println("Colonist with name \"" + mainColonist.Name + "\" created.")
	} else {
		mainColonist = g.CreateColonist()
	}
	return
}

type QuestionAsker struct {
	ui.InputOutput
}

func (q *QuestionAsker) Ask(question string) (answer string) {
	q.Println(question)
	answer, _ = q.Read()
	answer = strings.ToLower(answer)
	return
}

func (q *QuestionAsker) AskYesNo(question string) (yesno bool) {
	q.Println(question)
	answer, _ := q.Read()
	answer = strings.ToLower(answer)

	if (answer[0] == 'y') {
		yesno = true
	} else {
		yesno = false
	}
	return
}

type BaseFactory struct{
	QuestionAsker
}

func (b *BaseFactory) CreateBase() (base *world.Base) {
	name := b.Ask("What would you like to name your base?")
	base = &world.Base{Name: "Base" + name, Avatar: world.NewThing("Base " + name, "B", world.Position{2,2,2})}
	isAnswerYes := b.AskYesNo("Naming your base: \"" + base.Name + "\", is this correct? (y/n)")

	if (isAnswerYes) {
		b.Println("Base with name \"" + base.Name + "\" created.")
	} else {
		base = b.CreateBase()
	}
	return
}

func (g *GameManager) CreateBase() (base *world.Base){
	base = g.BaseFactory.CreateBase() // TODO Create Base Factory
	return
}


func  (g *GameManager) SetUpDirectionMaps() (err error) {
	letterToDirection := map[rune]Direction{'h':HERE,'n':NORTH,'e':EAST,'s':SOUTH,
		'w':WEST,'u':UP,'d':DOWN,'x':CANCEL,'p':SHOW_POSSIBILITIES}
	g.LetterToDirection = letterToDirection

	directionToString := map[Direction]string{HERE : "Here", NORTH : "North", EAST : "East",
		SOUTH : "South", WEST : "West", UP : "Up", DOWN : "Down"}
	g.DirectionToString = directionToString

	err = nil
	return
}


type Direction int

const (
	HERE Direction = iota
	NORTH
	EAST
	SOUTH
	WEST
	UP
	DOWN
	SHOW_POSSIBILITIES
	CANCEL
	RETRY
)

func (g *GameManager) Look(args []string) (err error) {
	w := g.World
	var dir Direction

	if (len(args) <= 1) {
		for dir = RETRY ; (dir == RETRY); {
			g.Println("In which direction would you like to look? ('c' to cancel or 'p' for possible directions)")
			dirString, _ := g.Read()
			dirString = strings.ToLower(dirString)
			if (strings.HasPrefix(dirString, "look")) {
				dirString = strings.TrimSpace(dirString[4:])
			}
			dir = g.GetDirection(dirString)
			if (dir == RETRY) {
				g.Println("The possible directions are H, N, E, S, W, U, D")
			}
		}
	} else {
		dirString := args[1]
		dir = g.GetDirection(dirString)
	}
	var name1, name2 string
	var things []world.Thing
	var isClear bool

	var pos world.Position

	switch dir {
	case SHOW_POSSIBILITIES:
		g.Println("The possible directions are H, N, E, S, W, U, D")
		return
	case CANCEL:
		g.Println("Cancelled looking")
		return
	case HERE: pos = world.Position{0,0,0}
	case NORTH: pos = world.Position{0,1,0}
	case EAST: pos = world.Position{1,0,0}
	case SOUTH: pos = world.Position{0,-1,0}
	case WEST: pos = world.Position{-1,0,0}
	case UP: pos = world.Position{0,0,1}
	case DOWN: pos = world.Position{0,0,-1}
	default:
		g.Println("The possible directions are H, N, E, S, W, U, D,")
		return
	}

	name1, name2, things, isClear, _ = w.GetNamesOfTerrainsAndObjects(pos)

	dirName := g.DirectionToString[dir]

	var middleString string
	var endString string
	if (isClear) {
		middleString =  " above " + name2
	}

	if (things != nil) && (len(things) > 0) {
		endString = " with: "
		for i, thing := range things {
			endString += thing.Name()
			if (i < len(things) - 2) {
				endString += ","
			} else if (i == len(things) - 2) {
				endString += ", and "
			}
		}
	}

	g.Println(dirName  + " is " + name1 + middleString + endString + ".")
	err = nil
	return
}

func (g *GameManager) Move(args []string) (err error) {
	w := g.World
	var dir Direction

	if (len(args) <= 1) {
		for dir = RETRY ; (dir == RETRY); {
			g.Println("In which direction would you like to move? ('c' to cancel or 'p' for possible directions)")
			dirString, e := g.Read()
			if (e != nil) {
				return e
			}
			dirString = strings.ToLower(dirString)
			if (strings.HasPrefix(dirString, "move")) {
				dirString = strings.TrimSpace(dirString[4:])
			}
			dir = g.GetDirection(dirString)
			if (dir == RETRY) {
				g.Println("The possible directions are H, N, E, S, W, U, D")
			}
		}
	} else {
		dirString := args[1]
		dir = g.GetDirection(dirString)
	}
	var name1, name2 string
	var things []world.Thing
	var isClear bool

	var pos world.Position

	switch dir {
	case SHOW_POSSIBILITIES:
		g.Println("The possible directions are H, N, E, S, W, U, D")
		return
	case CANCEL:
		g.Println("Cancelled moving")
		return
	case HERE: pos = world.Position{0,0,0}
	case NORTH: pos = world.Position{0,1,0}
	case EAST: pos = world.Position{1,0,0}
	case SOUTH: pos = world.Position{0,-1,0}
	case WEST: pos = world.Position{-1,0,0}
	case UP: pos = world.Position{0,0,1}
	case DOWN: pos = world.Position{0,0,-1}
	default:
		g.Println("The possible directions are H, N, E, S, W, U, D,")
		return
	}

	name1, name2, things, isClear, _ = w.GetNamesOfTerrainsAndObjects(pos)

	moveDirName := g.DirectionToString[dir]
	dirName := "Here"

	var middleString string
	var endString string
	if (isClear) {
		middleString =  " above " + name2
	}

	if (things != nil) && (len(things) > 0) {
		endString = " with: "
		for i, thing := range things {
			endString += thing.Name()
			if (i < len(things) - 2) {
				endString += ","
			} else if (i == len(things) - 2) {
				endString += ", and "
			}
		}
	}

	g.Println("Moved " + moveDirName + ".")
	g.Println(dirName  + " is " + name1 + middleString + endString + ".")

	err = g.World.ShiftColonist(pos)
	return
}

func  (g *GameManager) GetDirection(dirString string) (dir Direction) {
	dirChar := dirString[0]
	dir, isThere := g.LetterToDirection[rune(dirChar)]

	if !isThere {
		dir = RETRY
	}
	return
}