package world

import (
	"github.com/subspace-engine/subspace/util"
	"github.com/subspace-engine/subspace/world/model"
)

type BasicTile struct {
	model.Thing
	position util.Vec3
}

func MakeBasicTile(tileObject model.Thing) Tile {
	return &BasicTile{tileObject, util.Vec3{0, 0, 0}}
}

func (tile BasicTile) Position() util.Vec3 {
	space, ok := tile.Location().(Space)
	if ok {
		return tile.position.Mul(space.TileSize())
	}
	return tile.position
}

func (tile BasicTile) TilePosition() util.Vec3 {
	return tile.position
}

func (tile *BasicTile) SetPosition(pos util.Vec3) {
	tile.position = pos
}

var nothing = model.MakePassableThing("Nothing", "Nothing", false)

type Tiles [][][]Tile
type Things [][][][]model.Thing

type BasicSpace struct {
	tiles       Tiles
	things      Things
	tileSize    float64
	thingMul    float64
	players     []model.Player
	model.Thing // to conform to the Thing interface, allows us to give spaces names and descriptions
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
	space := BasicSpace{MakeTiles(width, height, depth, tile), MakeThings(width/int(thingMul), height/int(thingMul), depth/int(thingMul)), size, thingMul, make([]model.Player, 0, 0), model.MakeBasicThing("World", "The world")}
	return &space
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

func (space BasicSpace) TileSize() float64 {
	return space.tileSize
}

func (self BasicSpace) Move(thing model.Thing, pos util.Vec3) bool {
	tilepos := thing.Position().Div(self.tileSize)
	newpos := thing.Position().Add(pos).Div(self.tileSize)
	if self.Encloses(newpos) {
		tile := self.tiles[int(newpos.X)][int(newpos.Y)][int(newpos.Z)]
		if tile.Passable() {

			self.shiftThing(thing, tilepos, newpos)

			thing.SetPosition(thing.Position().Add(pos))
			thing.SetLocation(tile)
			go thing.Act(model.Action{thing, "move", tile, nil})
			self.calculateEncounters(tile, thing, newpos)
			return true
		}
		go thing.Act(model.Action{thing, "bump", tile, nil})
		return false
	}
	return false
}

func (self *BasicTile) Move(thing model.Thing, pos util.Vec3) bool {
	return self.Location().Move(thing, pos)
}

func (self *BasicSpace) shiftThing(thing model.Thing, pos util.Vec3, newpos util.Vec3) {
	thingMul := self.thingMul
	if !pos.Div(thingMul).Equals(newpos.Div(thingMul)) {

		self.things.remove(pos.Div(thingMul), thing)
		self.things.add(newpos.Div(thingMul), thing)
	}
}

func (self BasicSpace) GetTile(thing model.Thing) Tile {
	x := int(thing.Position().X / self.tileSize)
	y := int(thing.Position().Y / self.tileSize)
	z := int(thing.Position().Z / self.tileSize)
	return self.tiles[x][y][z]
}

func (self *BasicSpace) SetTile(x int, y int, z int, tile Tile) {
	self.tiles[x][y][z] = tile
	tile.SetPosition(util.Vec3{float64(x), float64(y), float64(z)})
	tile.SetLocation(self)
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

func (self *BasicSpace) Add(thing model.Thing) {
	tilepos := self.Position().Div(self.tileSize)
	if self.Encloses(tilepos) {
		self.things.add(tilepos.Div(self.thingMul), thing)
		thing.SetLocation(self.TileAt(int(tilepos.X), int(tilepos.Y), int(tilepos.Z)))
		player, ok := thing.(model.Player)
		if ok {
			self.addPlayer(player)
		}
	}
}

func (self *BasicSpace) addPlayer(player model.Player) {
	self.players = append(self.players, player)
}

func (self *BasicSpace) removePlayer(player model.Player) {
	for i, val := range self.players {
		if val == player {
			self.players[i] = self.players[len(self.players)-1]
			self.players = self.players[:len(self.players)-1]
		}
	}
}

func (self *BasicSpace) thingsOnTile(x int, y int, z int) []model.Thing {
	children := make([]model.Thing, 0, 0)
	tilePosition := util.Vec3{float64(x), float64(y), float64(z)}
	thingLocation := tilePosition.Div(self.thingMul)
	for _, val := range self.things[int(thingLocation.X)][int(thingLocation.Y)][int(thingLocation.Z)] {
		if val.Position().Div(self.tileSize).Equals(tilePosition) {
			children = append(children, val)
		}
	}
	return children
}

func (tile BasicTile) Children() []model.Thing {
	space, err := tile.Location().(*BasicSpace)
	if !err {
		return make([]model.Thing, 0, 0)
	}
	pos := tile.TilePosition()
	return space.thingsOnTile(int(pos.X), int(pos.Y), int(pos.Z))
}

func (space BasicSpace) IsRoot() bool {
	return true
}

func (self *BasicSpace) calculateEncounters(tile Tile, thing model.Thing, pos util.Vec3) {
	children := tile.Children()

	for _, val := range children {
		if val != thing {
			go thing.Act(model.Action{thing, "encounter", val, nil})
			go val.Act(model.Action{tile, "arrived", thing, nil})
		}
	}
}

func (space *BasicSpace) Remove(thing model.Thing) {
	pos := thing.Position().Div(space.tileSize).Div(space.thingMul)
	space.things.remove(pos, thing)
	thing.SetLocation(nil)
	player, ok := thing.(model.Player)
	if ok {
		space.removePlayer(player)
	}
}

func (self *BasicSpace) Say(text string) {
	for _, val := range self.players {
		val.Print(text)
	}
}
