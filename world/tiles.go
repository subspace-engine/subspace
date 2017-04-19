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

type Tile interface {
	IsPassable() bool
}

type BasicTile struct {
	Type TileType
}	

func MakeBasicTile(tileType TileType) Tile {
	return BasicTile {tileType}
}

func (self BasicTile)IsPassable() bool {
	return self.Type!=Wall&&self.Type!=Nothing
}

type Tiles [][][]Tile
type Movers [][][][]model.Mover

type Space struct {
	tiles Tiles
	movers Movers
	TileSize float64
	moverMul int
}

func MakeTiles(width int, height int, depth int, tile Tile) Tiles {
	tiles := make([][][]Tile, width)
	for i:=0; i< width; i++ {
		tiles[i] = make([][]Tile, height)
		for j:=0; j<height; j++ {
			tiles[i][j] = make([]Tile, depth)
			for k:=0; k < depth; k++ {
				tiles[i][j][k] = tile
			}
		}
	}
	return tiles
}

func MakeMovers(width int, height int, depth int) Movers {
	movers := make([][][][]model.Mover, width)
	for i:=0; i < width; i++ {
		movers[i] = make([][][]model.Mover, height)
		for j:=0; j < height; j++ {
			movers[i][j] = make([][]model.Mover, depth)
			for k:=0; k< depth; k++ {
				movers[i][j][k] = make([]model.Mover, 0)
			}
		}
	}
	return movers
}

func MakeSpace(width int, height int, depth int, size float64, moverMul int, tile Tile) Space {
	space := Space{MakeTiles(width, height, depth, tile), MakeMovers(width/moverMul, height/moverMul, depth/moverMul), size, moverMul}
	return space
}

func MakeDefaultSpace(width int, height int, depth int) Space {
	return MakeSpace(width, height, depth, 1.0, 10, MakeBasicTile(Nothing))
}

func (self Movers)remove(x int, y int, z int, mover model.Mover) {
	if len(self[x][y][z])==1 {
		self[x][y][z] = make([]model.Mover, 0)
	} else {
		for i:=0; i < len(self[x][y][z]); i++ {
			if self[x][y][z][i]==mover {
				self[x][y][z][i]= self[x][y][z][len(self[x][y][z])-1]
				break
			}
			self[x][y][z] = self[x][y][z][:len(self[x][y][z])-1]
		}
	}
}

func (self Movers)add(x int, y int, z int, mover model.Mover) {
	self[x][y][z] = append(self[x][y][z], mover)
}

func (self *Space)Move(mover model.Mover, x float64, y float64, z float64) int {
	tx := int(mover.X()/self.TileSize)
	ty := int(mover.Y()/self.TileSize)
	tz := int(mover.Z()/self.TileSize)
	nx :=int((mover.X()+x)/self.TileSize)
	ny:=int((mover.Y()+y)/self.TileSize)
	nz:=int((mover.Z()+z)/self.TileSize)
	if self.Encloses(nx,ny,nz) {
		if self.tiles[nx][ny][nz].IsPassable() {
			self.shiftMover(mover, tx, ty, tz, nx, ny, nz)
			mover.SetX(mover.X()+x)
			mover.SetY(mover.Y()+y)
			mover.SetZ(mover.Z()+z)
			return 0
		}
		return 1
	}
	return -1
}

func (self *Space)shiftMover(mover model.Mover, x int, y int, z int, newX int, newY int, newZ int) {
	moverMul := self.moverMul
	if x/moverMul!=newX/moverMul|| y/moverMul!=newY/moverMul || z/moverMul!=newZ/moverMul {
		self.movers.remove(x/moverMul, y/moverMul, z/moverMul, mover)
		self.movers.add(newX/moverMul, newY/moverMul, newZ/moverMul, mover)
	}
}

func (self Space)GetTile(mover model.Mover) Tile {
	return self.tiles[int(mover.X()/self.TileSize)][int(mover.Y()/self.TileSize)][int(mover.Z()/self.TileSize)]
}

func (self Space)SetTile(x int, y int, z int, tile Tile) {
	self.tiles[x][y][z]=tile
}

func (self Space)Encloses(x int, y int, z int) bool {
	return x>=0&&x<len(self.tiles)&&y>=0&&y<len(self.tiles[x])&&z>=0&&z<len(self.tiles[x][y])
}

func (self Space)Add(x int, y int, z int, mover model.Mover) {
	if self.Encloses(x,y,z) {
		self.movers.add(x/self.moverMul,y/self.moverMul,z/self.moverMul, mover)
	}
}