package world

// World, a type intended to hold all the objects of a game and how they relate to each other.

// UnsetDimention indicates when a world does not use a particular dimention, e.g. card games will use this for X, Y and ZSize, while side-scrolers will only use it for the ZSize
const UnsetDimention = 0.0

type TerrainType uint

func NewTerrain(x, y, z int) [][][]TerrainType {
	terrain := make([][][]TerrainType, x, x)
	for i := range terrain {
		terrain[i] = make([][]TerrainType, y, y)
		for j := range terrain[i] {
			terrain[i][j] = make([]TerrainType, z, z)
		}
	}
	return terrain
}

func NewTiles(x, y, z, sf int) [][][]ObjectTile {
	o1 := scaleSize(x, sf)
	o2 := scaleSize(y, sf)
	o3 := scaleSize(z, sf)
	tiles := make([][][]ObjectTile, o1, o1)
	for i := range tiles {
		tiles[i] = make([][]ObjectTile, o2, o2)
		for j := 0; j < o3; j++ {
			tiles[i][j] = make([]ObjectTile, o3, o3)
			for k := range tiles[i][j] {
				tiles[i][j][k] = make([]Tile, 0, 0)
			}
		}
	}
	return tiles
}

type Tile struct {
	xOffset uint
	yOffset uint
	zOffset uint
	// todo: hold a base object
}

type ObjectTile []Tile

type World struct {
	XSize float64
	YSize float64
	ZSize float64

	//below types can be optimised later to not be multidimentional arrays
	// Terrain represents any cube of world space that requires unique feature types, setting this to nil disables associated logic
	Terrain     [][][]TerrainType
	ObjectTiles [][][]ObjectTile
}

func NewWorld(x, y, z int, hasTerrain bool, sf int) *World {
	var tiles [][][]ObjectTile
	var terrain [][][]TerrainType
	if x != 0 || y != 0 || z != 0 {
		if hasTerrain {
			terrain = NewTerrain(x, y, z)
		}
		if sf != 0 {
			tiles = NewTiles(x, y, z, sf)
		}
	}
	world := &World{XSize: float64(x), YSize: float64(y), ZSize: float64(z), Terrain: terrain, ObjectTiles: tiles}
	return world
}

func scaleSize(orig, sf int) int {
	if orig < sf {
		return 1
	}
	size := orig / sf
	if orig%sf != 0 {
		size++
	}
	return size
}
