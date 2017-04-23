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
	Type() int
	IsPassable() bool
	String() string
}

type Space interface {
	Move(mover model.Mover, x float64, y float64, z float64) int
	GetTile(mover model.Mover) Tile
	SetTile(x int, y int, z int, tile Tile)
	Encloses(x int, y int, z int) bool
	Add(x int, y int, z int, mover model.Mover)
}
