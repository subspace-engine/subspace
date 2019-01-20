package world

import "github.com/subspace-engine/subspace/world/model"

type BasicTile struct {
	tileObject model.TileObject
}

func (self BasicTile)TileObject() model.TileObject {
	return self.tileObject
}

func MakeBasicTile(tileObject model.TileObject) Tile {
	return BasicTile{tileObject}
}

func (self BasicTile) IsPassable() bool {
	return self.TileObject().Passable()
}

func (tile BasicTile)Type() int {
	return tile.TileObject().Type()
}

func (tile BasicTile)String() string {
	return tile.TileObject().Name()
}

var nothing = model.MakePassableThing("Nothing", "Nothing", false)
type Tiles [][][]Tile
type Things [][][][]*model.Thing

type BasicSpace struct {
	tiles    Tiles
	things   Things
	TileSize float64
	thingMul int
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
	things := make([][][][]*model.Thing, width)
	for i := 0; i < width; i++ {
		things[i] = make([][][]*model.Thing, height)
		for j := 0; j < height; j++ {
			things[i][j] = make([][]*model.Thing, depth)
			for k := 0; k < depth; k++ {
				things[i][j][k] = make([]*model.Thing, 0)
			}
		}
	}
	return things
}

func MakeBasicSpace(width int, height int, depth int, size float64, thingMul int, tile Tile) Space {
	space := BasicSpace{MakeTiles(width, height, depth, tile), MakeThings(width/thingMul, height/thingMul, depth/thingMul), size, thingMul}
	return space
}

func MakeDefaultSpace(width int, height int, depth int) Space {
	return MakeBasicSpace(width, height, depth, 1.0, 10, MakeBasicTile(nothing))
}

func (self Things) remove(x int, y int, z int, thing *model.Thing) {
	if len(self[x][y][z]) == 1 {
		self[x][y][z] = make([]*model.Thing, 0)
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

func (self Things) add(x int, y int, z int, thing *model.Thing) {
	self[x][y][z] = append(self[x][y][z], thing)
}

func (self BasicSpace) Move(thing *model.Thing, x float64, y float64, z float64) int {
	tx := int(thing.X() / self.TileSize)
	ty := int(thing.Y() / self.TileSize)
	tz := int(thing.Z() / self.TileSize)
	nx := int((thing.X() + x) / self.TileSize)
	ny := int((thing.Y() + y) / self.TileSize)
	nz := int((thing.Z() + z) / self.TileSize)
	if self.Encloses(nx, ny, nz) {
		if self.tiles[nx][ny][nz].IsPassable() {
			self.shiftThing(thing, tx, ty, tz, nx, ny, nz)
			thing.SetX(thing.X() + x)
			thing.SetY(thing.Y() + y)
			thing.SetZ(thing.Z() + z)
						return 0
		}
		tile :=self.tiles[nx][ny][nz]
		thing.Act(model.Action{thing, "bump", tile.TileObject().(*model.Thing), nil})
		tile.TileObject().Act(model.Action{tile.TileObject().(*model.Thing), "bump", thing, nil})
		return 1
	}
	return -1
}

func (self *BasicSpace) shiftThing(thing *model.Thing, x int, y int, z int, newX int, newY int, newZ int) {
	thingMul := self.thingMul
	if x/thingMul != newX/thingMul || y/thingMul != newY/thingMul || z/thingMul != newZ/thingMul {
		self.things.remove(x/thingMul, y/thingMul, z/thingMul, thing)
		self.things.add(newX/thingMul, newY/thingMul, newZ/thingMul, thing)
	}
}

func (self BasicSpace) GetTile(thing *model.Thing) Tile {
	return self.tiles[int(thing.X()/self.TileSize)][int(thing.Y()/self.TileSize)][int(thing.Z()/self.TileSize)]
}

func (self BasicSpace) SetTile(x int, y int, z int, tile Tile) {
	self.tiles[x][y][z] = tile
}

func (self BasicSpace) Encloses(x int, y int, z int) bool {
	return x >= 0 && x < len(self.tiles) && y >= 0 && y < len(self.tiles[x]) && z >= 0 && z < len(self.tiles[x][y])
}

func (self BasicSpace) Add(x int, y int, z int, thing *model.Thing) {
	if self.Encloses(x, y, z) {
		self.things.add(x/self.thingMul, y/self.thingMul, z/self.thingMul, thing)
	}
}
