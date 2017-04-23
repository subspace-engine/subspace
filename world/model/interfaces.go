package model

type FPoint struct {
	X float64
	Y float64
	Z float64
}

type Mover interface {
	Pos() FPoint
	SetPos(pos FPoint)
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
	GetTile(mover Mover) Tile
	SetTile(pos Point, tile Tile)
	Encloses(tilePos Point, objPos FPoint) bool
	AddObject(pos FPoint, mover Mover)
	MoveObject(mover Mover, pos FPoint) int
}

type TileType int

type Tile interface {
	IsPassable() bool
	GetType() TileType
}

