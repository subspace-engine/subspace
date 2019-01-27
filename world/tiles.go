package world

import (
	"github.com/subspace-engine/subspace/util"
	"github.com/subspace-engine/subspace/world/model"
)

type BasicTile struct {
	model.Thing
}

func MakeBasicTile(tileObject model.Thing) Tile {
	return BasicTile{tileObject}
}

var nothing = model.MakePassableThing("Nothing", "Nothing", false)

type Tiles [][][]Tile
type Things [][][][]model.Thing

type BasicSpace struct {
	tiles    Tiles
	things   Things
	TileSize float64
	thingMul float64
}

func MakeTiles(width int, height int, depth int, tile Tile) Tiles {
	tiles := make([][][]Tile, width)
	for i := 0; i < width; i++ {
		tiles[i] = make([][]Tile, height)
		for j := 0; j < height; j++ {
			tiles[i][j] = make([]Tile, depth)
			for k := 0; k < depth; k++ {
				tiles[i][j][k] = tile
			}
		}
	}
	return tiles
}

func MakeThings(width int, height int, depth int) Things {
	things := make([][][][]model.Thing, width)
	for i := 0; i < width; i++ {
		things[i] = make([][][]model.Thing, height)
		for j := 0; j < height; j++ {
			things[i][j] = make([][]model.Thing, depth)
			for k := 0; k < depth; k++ {
				things[i][j][k] = make([]model.Thing, 0)
			}
		}
	}
	return things
}

func MakeBasicSpace(width int, height int, depth int, size float64, thingMul float64, tile Tile) Space {
	space := BasicSpace{MakeTiles(width, height, depth, tile), MakeThings(width/int(thingMul), height/int(thingMul), depth/int(thingMul)), size, thingMul}
	return space
}

func MakeDefaultSpace(width int, height int, depth int) Space {
	return MakeBasicSpace(width, height, depth, 1.0, 10, MakeBasicTile(nothing))
}

func (self Things) remove(pos util.Vec3, thing model.Thing) {
	x := int(pos.X)
	y := int(pos.Y)
	z := int(pos.Z)
	if len(self[x][y][z]) == 1 {
		self[x][y][z] = make([]model.Thing, 0)
	} else {
		for i := 0; i < len(self[x][y][z]); i++ {
			if self[x][y][z][i] == thing {
				self[x][y][z][i] = self[x][y][z][len(self[x][y][z])-1]
				break
			}
			self[x][y][z] = self[x][y][z][:len(self[x][y][z])-1]
		}
	}
}

func (self Things) add(pos util.Vec3, thing model.Thing) {
	x := int(pos.X)
	y := int(pos.Y)
	z := int(pos.Z)
	self[x][y][z] = append(self[x][y][z], thing)
}

func (self BasicSpace) Move(thing model.Thing, pos util.Vec3) int {
	tilepos := thing.Position().Div(self.TileSize)
	newpos := thing.Position().Add(pos).Div(self.TileSize)
	if self.Encloses(newpos) {
		tile := self.tiles[int(newpos.X)][int(newpos.Y)][int(newpos.Z)]
		if tile.Passable() {
			self.shiftThing(thing, tilepos, newpos)
			thing.SetPosition(thing.Position().Add(pos))
			go tile.Act(model.Action{tile, "step", thing, nil})
			go thing.Act(model.Action{thing, "step", tile, nil})
			return 0
		}
		go thing.Act(model.Action{thing, "bump", tile, nil})
		go tile.Act(model.Action{tile, "bump", thing, nil})
		return 1
	}
	return -1
}

func (self *BasicSpace) shiftThing(thing model.Thing, pos util.Vec3, newpos util.Vec3) {
	thingMul := self.thingMul
	if !pos.Div(thingMul).Equals(newpos.Div(thingMul)) {
		self.things.remove(pos.Div(thingMul), thing)
		self.things.add(newpos.Div(thingMul), thing)
	}
}

func (self BasicSpace) GetTile(thing model.Thing) Tile {
	x := int(thing.Position().X / self.TileSize)
	y := int(thing.Position().Y / self.TileSize)
	z := int(thing.Position().Z / self.TileSize)
	return self.tiles[x][y][z]
}

func (self BasicSpace) SetTile(x int, y int, z int, tile Tile) {
	self.tiles[x][y][z] = tile
}

func (self BasicSpace) TileAt(x int, y int, z int) Tile {
	return self.tiles[x][y][z]
}

func (self BasicSpace) Encloses(pos util.Vec3) bool {
	x := int(pos.X)
	y := int(pos.Y)
	z := int(pos.Z)
	return x >= 0 && x < len(self.tiles) && y >= 0 && y < len(self.tiles[x]) && z >= 0 && z < len(self.tiles[x][y])
}

func (self BasicSpace) Add(thing model.Thing) {
	if self.Encloses(thing.Position()) {
		self.things.add(thing.Position().Mul(self.thingMul), thing)
	}
}
