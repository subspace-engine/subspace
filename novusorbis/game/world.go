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
	Things       *MapThingStore
	Cursor       Position
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
	Avatar Thing
}

func (inPos *Position) RelativePosition(x int, y int, z int) (outPos Position, err error) {
	outPos = Position{inPos.x + x, inPos.y + y, inPos.z + z}
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
			symbol, _ := w.GetSymbolOfWorldAt(Position{x,y,z})
			byteWorld.WriteString(symbol)
		}
		byteWorld.WriteString("\n")
	}

	drawnWorld = byteWorld.String()
	return
}

func (w *World) GetSymbolOfWorldAt(p Position) (worldChar string, err error) {
	terrain := w.Terrain
	things, _ := w.Things.AtPosition(p)
	if len(things) > 0 {
		err = nil
		worldChar = things[0].Symbol()
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

	c := Position{mid,mid,mid}
	w.Cursor = c

	w.GenerateTerrain()

	store := &MapThingStore{}
	store.Initialize()
	w.Things = store

	w.MainColonist = g.CreateColonist()
	// g.CreateBase() // TODO return base

	pos := Position{mid,mid,mid}
	obj := g.World.MainColonist.Avatar
	obj.SetPosition(pos)
	store.AddObjectAt(obj, pos)

	err = nil
	return
}

func (g *GameManager) CreateColonist() (mainColonist *Colonist) {
	out := g.Out
	in := g.In

	out.Println("What would you like to name your first colonist?")
	name := in.Read()
	colonist := &Colonist{Avatar:&BasicThing{name : "you", symbol : "@"}, Name : name}
	mainColonist = colonist
	out.Println("Creating a colonist with the name: \"" + colonist.Name + "\", is this correct? (y/n)")
	answer := strings.ToLower(in.Read())
	if (len(answer) > 0 && answer[0] == 'y') {
		out.Println("Colonist with name \"" + colonist.Name + "\" created.")
	} else {
		mainColonist = g.CreateColonist()
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
	var things []Thing
	var isClear bool

	var pos Position

	switch dir {
	case SHOW_POSSIBILITIES:
		out.Println("The possible directions are H, N, E, S, W, U, D")
		return
	case CANCEL:
		out.Println("Cancelled looking")
		return
	case HERE: pos = Position{0,0,0}
	case NORTH: pos = Position{0,1,0}
	case EAST: pos = Position{1,0,0}
	case SOUTH: pos = Position{0,-1,0}
	case WEST: pos = Position{-1,0,0}
	case UP: pos = Position{0,0,1}
	case DOWN: pos = Position{0,0,-1}
	default:
		out.Println("The possible directions are H, N, E, S, W, U, D,")
		return
	}

	name1, name2, things, isClear, _ = world.GetNamesOfTerrainsAndObjects(pos)

	dirName := g.DirectionToString[dir]

	var middleString string
	var endString string
	if (isClear) {
		middleString =  " above " + name2
	}

	if (things != nil) && (len(things) > 0) {
		endString = " with: " + things[0].Name()
	}

	out.Println(dirName  + " is " + name1 + middleString + endString + ".")
	err = nil
	return
}

func (w *World) GetNamesOfTerrainsAndObjects(relPos Position) (name1 string,
		name2 string,
		things []Thing,
		isClear bool,
		err error) {
	pos, _ := w.Cursor.RelativePosition(relPos.x, relPos.y, relPos.z)
	terrainType, _ := w.Terrain.GetTerrainTypeAt(pos)
	name1, _ = w.Terrain.GetNameOfTerrainAt(pos)
	isClear = (w.Terrain.TerrainToOpacity[terrainType] == CLEAR)

	things, _ = w.Things.AtPosition(pos)

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