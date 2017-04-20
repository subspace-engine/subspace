package game

import (
	"strings"
	"bytes"
)

type World struct {
	MainColonist *Colonist
	MainBase *Base
	Size int
	Terrain *Terrain
	Objects *ObjectStore
	Cursor *Position
}

type ObjectStore struct {
	Objects map[Position]*GameObject
}

type GameObject struct {
	Name string
}

type Position struct {
	x int
	y int
	z int
}

type Base struct {
	Name string
}

type Colonist struct {
	Name string
}

func (store *ObjectStore) Initialize() {
	store.Objects = make(map[Position]*GameObject)
}

func (store *ObjectStore) AtPosition(p *Position) (object *GameObject, isFound bool, err error) {
	object, isFound = store.Objects[*p]
	err = nil
	return
}


func (store *ObjectStore) AddObjectAt(obj *GameObject, p *Position) (err error) {
	store.Objects[*p] = obj
	err = nil
	return
}

func (inPos *Position) RelativePosition(x int, y int, z int) (outPos *Position, err error) {
	outPos = &Position{inPos.x + x, inPos.y + y, inPos.z + z}
	err = nil
	return
}

func (w *World) DrawnWorldAtZ(z int) (drawnWorld string, err error){
	size := w.Size
	err = nil
	var byteWorld bytes.Buffer

	if z >= size || z < 0 {
		drawnWorld = "z value out of range\n"
		return
	}

	for y:=0; y < size; y++ {
		for x := 0; x < size; x++ {
			symbol, _ := w.GetSymbolOfWorldAt(&Position{x,y,z})
			byteWorld.WriteString(symbol)
		}
		byteWorld.WriteString("\n")
	}

	drawnWorld = byteWorld.String()
	return
}

func (w *World) GetSymbolOfWorldAt(p *Position) (worldChar string, err error) {
	terrain := w.Terrain
	_, isFound, _ := w.Objects.AtPosition(p)
	if (isFound) {
		err = nil
		worldChar = "X"
	} else {
		worldChar, err = terrain.GetSymbolOfTerrainAt(p)
	}

	return
}


func (g *GameManager) CreateWorld() (err error) {
	w := &World{}
	w.Size = 5
	g.World = w
	mid := w.Size/2

	c := &Position{mid,mid,mid} // TODO Proper sizes
	w.Cursor = c

	store := &ObjectStore{}
	store.Initialize()
	w.Objects = store

	pos := &Position{mid,mid,mid}
	obj := &GameObject{"Random Object"}
	store.AddObjectAt(obj, pos)

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
		name1, _ := g.World.Terrain.GetNameOfTerrainAt(g.World.Cursor)
		pos, _ := g.World.Cursor.RelativePosition(0,0,-1)
		name2, _ := g.World.Terrain.GetNameOfTerrainAt(pos)
		out.Println("Here you see " + name1 + " above " + name2)
		obj, isFound, _ := g.World.Objects.AtPosition(g.World.Cursor)
		if (isFound) {
			out.Println("Here there is also " + obj.Name)
		}
	case NORTH:
		pos, _ := g.World.Cursor.RelativePosition(0,1,0)
		name, _ := g.World.Terrain.GetNameOfTerrainAt(pos)
		out.Println("North you see " + name)
	case EAST:
		pos, _ := g.World.Cursor.RelativePosition(1,0,0)
		name, _ := g.World.Terrain.GetNameOfTerrainAt(pos)
		out.Println("East you see " + name)
	case SOUTH:
		pos, _ := g.World.Cursor.RelativePosition(0,-1,0)
		name, _ := g.World.Terrain.GetNameOfTerrainAt(pos)
		out.Println("South you see " + name)
	case WEST:
		pos, _ := g.World.Cursor.RelativePosition(-1,0,0)
		name, _ := g.World.Terrain.GetNameOfTerrainAt(pos)
		out.Println("West you see " + name)
	case UP:
		pos, _ := g.World.Cursor.RelativePosition(0,0,1)
		name, _ := g.World.Terrain.GetNameOfTerrainAt(pos)
		out.Println("Up you see " + name)
	case DOWN:
		pos, _ := g.World.Cursor.RelativePosition(0,0,-1)
		name, _ := g.World.Terrain.GetNameOfTerrainAt(pos)
		out.Println("Down you see " + name)
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