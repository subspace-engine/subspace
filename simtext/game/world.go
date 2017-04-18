package game

import "strings"

type World struct {
	MainColonist *Colonist
	MainBase *Base
	Terrain *Terrain
}

type Base struct {
	Name string
}

type Colonist struct {
	Name string
}

func (g *GameManager) CreateWorld() (err error) {
	w := &World{}
	g.World = w
	// g.CreateColonist()
	// g.CreateBase()
	w.GenerateTerrain()

	err = nil
	return
}

func (g *GameManager) CreateColonist() {
	out := g.Out
	in := g.In
	world := g.World

	out.Println("What would you like to name your first colonist?")
	name := in.Read()
	colonist := &Colonist{Name : name}
	world.MainColonist = colonist
	out.Println("Creating a colonist with the name: \"" + colonist.Name + "\", is this correct? (y/n)")
	answer := strings.ToLower(in.Read())
	if (answer[0] == 'y') {
		out.Println("Colonist with name \"" + colonist.Name + "\" created.")
	} else {
		g.CreateColonist()
	}
	return
}

func (g *GameManager) CreateBase() {
	out := g.Out
	in := g.In
	out.Println("What would you like to name your base?")
	name := in.Read()
	world := g.World
	base := &Base{Name : name}
	world.MainBase = base
	out.Println("Naming your base: \"" + base.Name + "\", is this correct? (y/n)")
	answer := strings.ToLower(in.Read())
	if (answer[0] == 'y') {
		out.Println("Base with name \"" + base.Name + "\" created.")
	} else {
		g.CreateBase()
	}
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
	out := g.Out
	in := g.In
	var dir Direction

	if (len(args) <= 1) {
		for dir = RETRY ; (dir == RETRY); {
			out.Println("In which direction would you like to look? ('x' to cancel or 'p' for possible directions)")
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
	case HERE:
		out.Println("You looked Here!")
	case SOUTH:
		out.Println("You looked South!")
	case WEST:
		out.Println("You looked West!")
	case UP:
		out.Println("You looked Up!")
	case DOWN:
		out.Println("You looked Down!")
	case SHOW_POSSIBILITIES:
		out.Println("The possible directions are H, N, E, S, W, U, D, C")
	case CANCEL:
		out.Println("You cancelled looking")
	default:
		out.Println("I don't know which way you looked")
	}

	return nil
}

func  (g *GameManager) GetDirection(dirString string) (dir Direction) {

	dirChar := dirString[0]
	switch dirChar {
	case 'h':
		dir = HERE
	case 'n':
		dir = NORTH
	case 'e':
		dir = EAST
	case 's':
		dir = SOUTH
	case 'w':
		dir = WEST
	case 'u':
		dir = UP
	case 'd':
		dir = DOWN
	case 'x':
		dir = CANCEL
	case 'p':
		dir = SHOW_POSSIBILITIES
	default:
		dir = RETRY
	}
	return
}