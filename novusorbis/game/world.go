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
	Structures   *MapStructureStore
	Cursor       Position
}

type Position struct {
	x int
	y int
	z int
}

type Base struct {
	Name string
	Avatar Thing
}

type Colonist struct {
	Name string
	Avatar Thing
}

func (w *World) ShiftThing(thing Thing, pos Position) (err error) {
	w.Things.ShiftObjBy(thing, pos)
	newPos, err := thing.Position().RelativePosition(pos.x, pos.y, pos.z)
	thing.SetPosition(newPos)
	err = nil
	return
}

func (w *World) ShiftColonist(pos Position) (err error) {
	c := w.MainColonist
	w.ShiftThing(c.Avatar, pos)
	err = w.ShiftCursor(pos)
	return
}

func (w *World) ShiftCursor(pos Position) (err error) {
	w.Cursor, err = w.Cursor.RelativePosition(pos.x, pos.y, pos.z)
	return err
}

func (inPos Position) RelativePosition(x int, y int, z int) (outPos Position, err error) {
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

	for y:=size-1; y >= 0; y-- {
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
	err = nil

	terrain := w.Terrain
	things, _ := w.Things.AtPosition(p)
	if len(things) == 0 {
		worldChar, err = terrain.GetSymbolOfTerrainAt(p)
		return
	}

	for _, thing := range things {
		if (thing == w.MainColonist.Avatar) {
			worldChar = thing.Symbol()
			return
		}
	}
	worldChar = things[0].Symbol()
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

	thingStore := &MapThingStore{}
	thingStore.Initialize()
	w.Things = thingStore

	structureStore := &MapStructureStore{}
	structureStore.Initialize()
	w.Structures = structureStore

	w.MainColonist = g.CreateDefaultColonist() // TODO
	w.MainBase = g.CreateDefaultBase() // TODO

	pos := Position{mid,mid,mid}
	mainAvatar := g.World.MainColonist.Avatar
	mainAvatar.SetPosition(pos)
	thingStore.AddObjectAt(mainAvatar, pos)

	baseAvatar := g.World.MainBase.Avatar
	baseAvatar.SetPosition(pos)
	thingStore.AddObjectAt(baseAvatar, pos)

	err = nil
	return
}

func (g *GameManager) CreateDefaultColonist() (mainColonist *Colonist) {
	mainColonist = &Colonist{Avatar:&BasicThing{name : "you", symbol : "@"}, Name : "Mark"}
	return
}

func (g *GameManager) CreateDefaultBase() (base *Base) {
	base = &Base{Avatar:&BasicThing{name : "Base Omicron", symbol : "b"}, Name : "Omicron"}
	return
}


func (g *GameManager) CreateColonist() (mainColonist *Colonist) {
	out := g.Out
	in := g.In

	out.Println("What would you like to name your first colonist?")
	name := in.Read()
	mainColonist = &Colonist{Avatar:&BasicThing{name : "you", symbol : "@"}, Name : name}
	out.Println("Creating a colonist with the name: \"" + mainColonist.Name + "\", is this correct? (y/n)")
	answer := strings.ToLower(in.Read())
	if (len(answer) > 0 && answer[0] == 'y') {
		out.Println("Colonist with name \"" + mainColonist.Name + "\" created.")
	} else {
		mainColonist = g.CreateColonist()
	}
	return
}

func (g *GameManager) CreateBase() (base *Base){
	out := g.Out
	in := g.In
	out.Println("What would you like to name your base?")
	name := in.Read()
	base = &Base{Avatar:&BasicThing{name : "Base " + name, symbol : "b"}, Name : name}
	out.Println("Naming your base: \"" + base.Name + "\", is this correct? (y/n)")
	answer := strings.ToLower(in.Read())
	if (answer[0] == 'y') {
		out.Println("Base with name \"" + base.Name + "\" created.")
	} else {
		base = g.CreateBase()
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

	out.Println(dirName  + " is " + name1 + middleString + endString + ".")
	err = nil
	return
}

func (g *GameManager) Move(args []string) (err error) {
	out := g.Out
	in := g.In
	world := g.World
	var dir Direction

	if (len(args) <= 1) {
		for dir = RETRY ; (dir == RETRY); {
			out.Println("In which direction would you like to move? ('c' to cancel or 'p' for possible directions)")
			dirString := strings.ToLower(in.Read())
			if (strings.HasPrefix(dirString, "move")) {
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
		out.Println("Cancelled moving")
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

	out.Println("Moved " + moveDirName + ".")
	out.Println(dirName  + " is " + name1 + middleString + endString + ".")

	err = g.World.ShiftColonist(pos)
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