package world

import (
	"github.com/subspace-engine/subspace/util"
	"github.com/subspace-engine/subspace/world/model"
)

type Tile interface {
	model.Thing
}

type Space interface {
	Move(model.Thing, util.Vec3) int
	GetTile(model.Thing) Tile
	SetTile(x int, y int, z int, tile Tile)
	Encloses(util.Vec3) bool
	TileAt(x int, y int, z int) Tile
	Add(model.Thing)
}
