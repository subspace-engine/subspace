package game

import (
	"strings"
	"bytes"
)

type World struct {
	MainColonist *Colonist
	MainBase     *Base
	Size         int
	Terrain      *Terrain
	Things       *ThingStore
	Cursor       *Position
}

type ThingStore struct {
	Things map[Position]*Thing
}

type Thing struct {
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

func (store *ThingStore) Initialize() {
	store.Things = make(map[Position]*Thing)
}

func (store *ThingStore) AtPosition(p *Position) (thing *Thing, isFound bool, err error) {
	thing, isFound = store.Things[*p]
	err = nil
	return
}

func (store *ThingStore) AddObjectAt(obj *Thing, p *Position) (err error) {
	store.Things[*p] = obj
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
	_, isFound, _ := w.Things.AtPosition(p)
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

	c := &Position{mid,mid,mid}
	w.Cursor = c

	store := &ThingStore{}
	store.Initialize()
	w.Things = store

	pos := &Position{mid,mid,mid}
	obj := &Thing{"Random Object"}
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

func (g *GameManager) Look(args []string) (err error) {
	out := g.Out
	in := g.In
	world := g.World
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
	var name1, name2 string
	var objects []*Thing
	var isClear bool

	switch dir {
	case HERE:
		name1, name2, objects, isClear, _ = world.GetNamesOfTerrainsAndObjects(Position{0,0,0})
	case NORTH:
		name1, name2, objects, isClear, _ = world.GetNamesOfTerrainsAndObjects(Position{0,1,0})
	case EAST:
		name1, name2, objects, isClear, _ = world.GetNamesOfTerrainsAndObjects(Position{1,0,0})
	case SOUTH:
		name1, name2, objects, isClear, _ = world.GetNamesOfTerrainsAndObjects(Position{0,-1,0})
	case WEST:
		name1, name2, objects, isClear, _ = world.GetNamesOfTerrainsAndObjects(Position{-1,0,0})
	case UP:
		name1, name2, objects, isClear, _ = world.GetNamesOfTerrainsAndObjects(Position{0,0,1})
	case DOWN:
		name1, name2, objects, isClear, _ = world.GetNamesOfTerrainsAndObjects(Position{0,0,-1})
	case SHOW_POSSIBILITIES:
		out.Println("The possible directions are H, N, E, S, W, U, D, C")
	case CANCEL:
		out.Println("You cancelled looking")
	default:
		out.Println("I don't know which way you looked")
	}

	dirName := g.DirectionToString[dir]

	var middleString string
	var endString string
	if (isClear) {
		middleString =  " above " + name2
	}

	if (objects != nil) {
		endString = " with " + objects[0].Name
	}

	out.Println(dirName  + " is " + name1 + middleString + endString + ".")
	err = nil
	return
}

func (w *World) GetNamesOfTerrainsAndObjects(relPos Position) (name1 string,
		name2 string,
		objects []*Thing,
		isClear bool,
		err error) {
	pos, _ := w.Cursor.RelativePosition(relPos.x, relPos.y, relPos.z)
	terrainType, _ := w.Terrain.GetTerrainTypeAt(pos)
	name1, _ = w.Terrain.GetNameOfTerrainAt(pos)
	isClear = (w.Terrain.TerrainToOpacity[terrainType] == CLEAR)

	obj, isFound, _ := w.Things.AtPosition(pos)

	if (isClear) {
		pos, _ = pos.RelativePosition(0, 0, -1)
		name2, _ = w.Terrain.GetNameOfTerrainAt(pos)
	}

	if (isFound) {
		objects = []*Thing{
			obj,
		}
	} else {
		objects = nil
	}

	err = nil
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