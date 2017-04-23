package game

import (
	"strings"
	"bytes"
	"github.com/subspace-engine/subspace/world/model"
	"fmt"
)

type World struct {
	MainColonist *Colonist
	MainBase     *Base
	Size         int
	Terrain      *Terrain
	Things       *ThingStore
	Cursor       *Point
}

type FPoint model.FPoint

type Thing struct {
	Name string
	Position FPoint
}

type ThingStore struct {
	Things []*Thing
}

func (thing *Thing) Pos() (p FPoint){
	p = thing.Position
	return
}

func (thing *Thing) SetPos(p FPoint){
	thing.Position = p
	return
}

type Point model.Point

type Base struct {
	Name string
}

type Colonist struct {
	Name string
}

func (store *ThingStore) Initialize() {
	store.Things = make([]*Thing,0,5)
}

func (world *World) AtPosition(p Point) (things []*Thing, err error) {
	things = make([]*Thing,0,5)
	fmt.Println("Getting at position ", p)

	for _, t := range world.Things.Things {
		fmt.Println("Thing ", t)
		if (world.Encloses(p, t.Pos())) {
			fmt.Println("Does enclose")
			things = append(things, t)
		}
	}

	err = nil
	return
}

func (world *World) Encloses(tilePos Point, objPos FPoint) bool {
	fmt.Println("Tilepos ", tilePos)
	fmt.Println("Objpos ", objPos)
	if (!(float64(tilePos.X) <= objPos.X && objPos.X <= float64(tilePos.X+1))) {
		return false
	}

	if (!(float64(tilePos.Y) <= objPos.Y && objPos.Y <= float64(tilePos.Y+1))) {
		return false
	}

	if (!(float64(tilePos.Z) <= objPos.Z  && objPos.Z <= float64(tilePos.Z+1))) {
		return false
	}
	return true
}



func (store *ThingStore) AddObjectAt(obj *Thing, p FPoint) (err error) {
	const DEFAULT_STORE_SIZE = 3
	store.Things = append(store.Things, obj)
	err = nil
	return
}

func (inPos Point) RelativePosition(x int, y int, z int) (outPos Point, err error) {
	outPos = Point{inPos.X + x, inPos.Y + y, inPos.Z + z}
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
			symbol, _ := w.GetSymbolOfWorldAt(Point{x,y,z})
			byteWorld.WriteString(symbol)
		}
		byteWorld.WriteString("\n")
	}

	drawnWorld = byteWorld.String()
	return
}

func (w *World) GetSymbolOfWorldAt(p Point) (worldChar string, err error) {
	terrain := w.Terrain
	things, _ := w.AtPosition(p)
	if len(things) > 0 {
		fmt.Println("Length of things was long")
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
	midFloat := float64(w.Size)/2
	midInt :=  w.Size/2

	c := &Point{midInt,midInt,midInt}
	w.Cursor = c

	store := &ThingStore{}
	store.Initialize()
	w.Things = store

	fmt.Println("MidFloat: ", midFloat)
	pos := FPoint{midFloat,midFloat,midFloat}
	obj := &Thing{"Random Object", pos}
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
			out.Println("In which direction would you like to look? ('c' to cancel or 'p' for possible directions)")
			dirString := strings.ToLower(in.Read())
			if (strings.HasPrefix(dirString, "look")) {
				dirString = strings.TrimSpace(dirString[4:])
			}
			dir = g.GetDirection(dirString)
			if (dir == RETRY) {
				out.Println("The possible directions are H, N, E, S, W, U, D")
			}
		}
	} else {
		dirString := args[1]
		dir = g.GetDirection(dirString)
	}
	var name1, name2 string
	var things []*Thing
	var isClear bool

	switch dir {
	case SHOW_POSSIBILITIES:
		out.Println("The possible directions are H, N, E, S, W, U, D")
		return
	case CANCEL:
		out.Println("Cancelled looking")
		return
	case HERE:
		name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(Point{0,0,0})
	case NORTH:
		name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(Point{0,1,0})
	case EAST:
		name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(Point{1,0,0})
	case SOUTH:
		name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(Point{0,-1,0})
	case WEST:
		name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(Point{-1,0,0})
	case UP:
		name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(Point{0,0,1})
	case DOWN:
		name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(Point{0,0,-1})
	default:
		out.Println("The possible directions are H, N, E, S, W, U, D,")
		return
	}

	dirName := g.DirectionToString[dir]

	var middleString string
	var endString string
	if (isClear) {
		middleString =  " above " + name2
	}

	if (things != nil) && (len(things) > 0) {
		endString = " with: " + things[0].Name
	}

	out.Println(dirName  + " is " + name1 + middleString + endString + ".")
	err = nil
	return
}

func (w *World) GetNamesOfTerrainsAndObjects(relPos Point) (name1 string,
		name2 string,
		things []*Thing,
		isClear bool,
		err error) {
	pos, _ := w.Cursor.RelativePosition(relPos.X, relPos.Y, relPos.Z)
	terrainType, _ := w.Terrain.GetTerrainTypeAt(pos)
	name1, _ = w.Terrain.GetNameOfTerrainAt(pos)
	isClear = (w.Terrain.TerrainToOpacity[terrainType] == CLEAR)

	things, _ = w.AtPosition(pos)

	if (isClear) {
		pos, _ = pos.RelativePosition(0, 0, -1)
		name2, _ = w.Terrain.GetNameOfTerrainAt(pos)
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