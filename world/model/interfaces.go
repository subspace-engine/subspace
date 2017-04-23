package model

type Mover interface {
	X() float64
	Y() float64
	Z() float64
	SetX(float64)
	SetY(float64)
	SetZ(float64)
}

type Namer interface {
	Name() string
	SetName(string)
}

type Describer interface {
	Description() string
	SetDescription(string)
}

type Point struct {
	X int
	Y int
	Z int
}

type Space interface {
	Move(mover Mover, x float64, y float64, z float64) int
	GetTile(mover Mover) Tile
	SetTile(pos Point, tile Tile)
	Encloses(pos Point) bool
	AddObject(pos Point, mover Mover)
}

type TileType int

type Tile interface {
	IsPassable() bool
	GetType() TileType
}

