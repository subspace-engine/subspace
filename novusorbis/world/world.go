package world

import (
	"bytes"
)

type World struct {
	MainColonist *Colonist
	MainBase     *Base
	Size         int
	Terrain      *Terrain
	Things       MoverStore
	Structures   MoverStore
	Cursor       Position
}

type Position struct {
	X int
	Y int
	Z int
}

type Base struct {
	Name   string
	Avatar Structure
}

type Colonist struct {
	Name   string
	Avatar Mover
	Inventory Container
}

func NewPosition(midpoint int) (pos Position) {
	pos = Position{X: midpoint, Y: midpoint, Z: midpoint}
	return
}

func NewDefaultColonist() (*Colonist){
	return &Colonist{"You", NewMover("Mark", "@" , Position{2,2,2}), NewContainer()}
}

func NewDefaultBase() (*Base){
	return &Base{Name : "Base Omicron", Avatar : NewStructure("Base", "B", Position{2,2,2})}
}

func (w *World) ShiftThing(thing Mover, pos Position) (err error) {
	w.Things.ShiftObjBy(thing, pos)
	newPos, err := thing.Position().RelativePosition(pos.X, pos.Y, pos.Z)
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
	w.Cursor, err = w.Cursor.RelativePosition(pos.X, pos.Y, pos.Z)
	return err
}

func (inPos Position) RelativePosition(x int, y int, z int) (outPos Position, err error) {
	outPos = Position{inPos.X + x, inPos.Y + y, inPos.Z + z}
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

func (w *World) GetNamesOfTerrainsAndObjects(relPos Position) (name1 string,
		name2 string,
		things []Mover,
		isClear bool,
		err error) {
	pos, _ := w.Cursor.RelativePosition(relPos.X, relPos.Y, relPos.Z)
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

