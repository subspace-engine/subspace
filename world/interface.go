package world

import "github.com/subspace-engine/subspace/world/model"

type Point interface {
	X() float64
	Y() float64
	Z() float64
	SetX(float64)
	SetY(float64)
	SetZ(float64)
}

type Tile interface {
	model.Thing
}

type Space interface {
	Move(model.Thing, float64, float64, float64) int
	GetTile(mover model.Thing) Tile
	SetTile(x int, y int, z int, tile Tile)
	Encloses(x int, y int, z int) bool
	TileAt(x int, y int, z int) Tile
	Add(x int, y int, z int, mover model.Thing)
}
