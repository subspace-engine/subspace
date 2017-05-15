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
	g.PrintLogo()
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
	g.CommandsMap["take"] = g.TakeObject
	g.CommandsMap["inv"] = g.ShowInventory
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
	return nil // TODO
	/*
	w := &world.World{}
	w.Size = 5
	g.World = w
	mid := w.Size/2

	c := world.Position{mid,mid,mid}
	w.Cursor = c

	w.GenerateTerrain()

	thingStore := &world.MapMoverStore{}
	thingStore.Initialize()
	w.Things = thingStore

	w.Structures = world.MapMoverStore{}

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
	*/
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

	mainColonist = &world.Colonist{Name: name, Avatar:world.NewMover("You", "@", world.Position{2,2,2})}
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
	base = &world.Base{Name: "Base" + name, Avatar: world.NewStructure("Base " + name, "B", world.Position{2,2,2})}
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

func (g *GameManager) TakeObject(args []string) (err error) {
	// TODO specify which object(s) to take
	// TODO differentiate between mobile and immobile objects
	colonist := g.World.MainColonist
	pos := colonist.Avatar.Position()
	thingsHere, err := g.World.Things.AtPosition(pos)
	if (len(thingsHere) <= 1) {
		g.Println("Nothing to take.")
		return
	}

	g.Println("Took: ")
	for index, thing := range thingsHere {
		if (thing == colonist.Avatar) {
			continue
		}
		g.Println(strconv.Itoa(index) + ". " + thing.Name())
		colonist.Inventory.AddObject(thing)
		g.World.Things.Remove(thing, pos)
	}
	return
}

func (g *GameManager) ShowInventory(args []string) (err error) {
	colonist := g.World.MainColonist
	inventory := colonist.Inventory.GetContents()
	if (len(inventory) == 0) {
		g.Println("Inventory empty")
		return
	}
	g.Println("Inventory: ")
	for index, thing := range inventory {
		g.Println(strconv.Itoa(index) + ". " + thing.Name())
	}
	return
}
