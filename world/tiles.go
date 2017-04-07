package world

import "github.com/subspace-engine/subspace/world/model"

type TileType int

const (
	Nothing TileType = iota
	Empty
	Wall
	Floor
	Ground
	LastType
)

func (self TileType)Text() string {
	var texts =[]string{"Nothing", "Empty", "Wall", "Floor", "Ground"}
	return texts[self]
}

type Tile struct {
	Type TileType
	Contents []model.Mover
}	

func MakeTile(tileType TileType) Tile {
	return Tile {tileType, make([]model.Mover, 0)}
}

func (self*Tile)IsPassable() bool {
	return self.Type!=Wall&&self.Type!=Nothing
}

func (self*Tile)Remove(mover model.Mover) {
	if len(self.Contents)<=1 {
		self.Contents = make([]model.Mover, 0)
	} else {
		for i:=0; i < len(self.Contents); i++ {
			if self.Contents[i]==mover {
				self.Contents[i]= self.Contents[len(self.Contents)-1]
				break
			}
			self.Contents = self.Contents[:len(self.Contents)-1]
		}
	}
}

func (self*Tile)Add(mover model.Mover) {
	self.Contents = append(self.Contents, mover)
}

type Tiles [][][]Tile

func MakeTiles(width int, height int, depth int) Tiles {
	tiles := make([][][]Tile, width)
	for i:=0; i< width; i++ {
		tiles[i] = make([][]Tile, height)
		for j:=0; j<height; j++ {
			tiles[i][j] = make([]Tile, depth)
			for k:=0; k < depth; k++ {
				tiles[i][j][k] = MakeTile(Nothing)
			}
		}
	}
	return tiles
}

func (self Tiles)Move(mover model.Mover, x float64, y float64, z float64) {
	nx :=int(mover.X()+x)
	ny:=int(mover.Y()+y)
	nz:=int(mover.Z()+z)
	if nx>=0&&nx<len(self)&&ny>=0&&ny<len(self[nx])&&nz>=0&&nz<len(self[nx][ny]) {
		if self[nx][ny][nz].IsPassable() {
			self[int(mover.X())][int(mover.Y())][int(mover.Z())].Remove(mover)
			mover.SetX(mover.X()+x)
			mover.SetY(mover.Y()+y)
			mover.SetZ(mover.Z()+z)
			self[nx][ny][nz].Add(mover)
		}
	}
}

func (self Tiles)GetTile(mover model.Mover) Tile {
	return self[int(mover.X())][int(mover.Y())][int(mover.Z())]
}

func (self Tiles)SetTile(x int, y int, z int, tile Tile) {
	self[x][y][z]=tile
}
