package world

// World, a type intended to hold all the objects of a game and how they relate to each other.

import (
	"fmt"
	"time"
)

// UnsetDimention indicates when a world does not use a particular dimention, e.g. card games will use this for X, Y and ZSize, while 2-dimentional worlds will only use it for the ZSize
const UnsetDimention = 0.0

type TerrainType uint8

const TerrainUnset = 0

func NewTerrain(x, y, z uint32) [][][]TerrainType {
	terrain := make([][][]TerrainType, x, x)
	for i := range terrain {
		terrain[i] = make([][]TerrainType, y, y)
		for j := range terrain[i] {
			terrain[i][j] = make([]TerrainType, z, z)
		}
	}
	return terrain
}

func NewTiles(x, y, z, sf uint32) [][][]*ObjectTile {
	o1 := scaleSize(x, sf)
	o2 := scaleSize(y, sf)
	o3 := scaleSize(z, sf)
	tiles := make([][][]*ObjectTile, o1, o1)
	for i := range tiles {
		tiles[i] = make([][]*ObjectTile, o2, o2)
		for j := 0; j < int(o3); j++ {
			tiles[i][j] = make([]*ObjectTile, o3, o3)
			for k := range tiles[i][j] {
				tiles[i][j][k] = &ObjectTile{X: i, Y: j, Z: k, Objects: make([]Tile, 0, 0)}
			}
		}
	}
	return tiles
}

type Tile struct {
	xOffset uint
	yOffset uint
	zOffset uint
	Object  interface{}
}

type ObjectTile struct {
	X, Y, Z int
	Objects []Tile
}

func (ot *ObjectTile) Add(t Tile) {
	ot.Objects = append(ot.Objects, t)
}

type Config struct {
	XSize uint32
	YSize uint32
	ZSize uint32
	ScaleFactor uint32

	// Update rate is the rate at wich the update loop is run
	UpdateRate  time.Time
	UseTerrain  bool
	UseTiles    bool
	// todo, more
}

type World struct {
	XSize uint32
	YSize uint32
	ZSize uint32
	SF    uint32 // scale factor of sectors to object tiles

	//below types can be optimised later to not be multidimentional arrays
	// Terrain represents any cube of world space that requires unique feature types, setting this to nil disables associated logic
	Terrain       [][][]TerrainType
	ObjectTiles   [][][]*ObjectTile
	GlobalObjects map[string]interface{}
}

func NewWorld(x, y, z uint32, hasTerrain bool, sf uint32) *World {
	var tiles [][][]*ObjectTile
	var terrain [][][]TerrainType
	if x != 0 || y != 0 || z != 0 {
		if hasTerrain {
			terrain = NewTerrain(x, y, z)
		}
		if sf != 0 {
			tiles = NewTiles(x, y, z, sf)
		}
	}
	return &World{XSize: x, YSize: y, ZSize: z, SF: sf, Terrain: terrain, ObjectTiles: tiles, GlobalObjects: make(map[string]interface{})}
}

func scaleSize(orig, sf uint32) uint32 {
	if orig < sf {
		return 1
	}
	size := orig / sf
	if orig%sf != 0 {
		size++
	}
	return size
}

func (w *World) AddObject(x, y, z int, obj interface{}) error {
	//todo, bounds check
	const offset = 0 // assume scale 1 for now, todo other sfs
	w.ObjectTiles[x][y][z].Add(Tile{offset, offset, offset, obj})
	return nil
}

type ActorFunc func(actor *Actor, args ...interface{})

type Actor struct {
	World     *World
	Actions   []ActorFunc
	X         int
	Y         int
	Z         int
	Objects   []*Tile
	Terrain   TerrainType
	CoordsSet bool
}

// pass actors by value to reuse them without changing the basis Actor
// I think this is what we usually want
// pass using a pointer to mutate
func (w *World) BuildActor(actions ...ActorFunc) Actor {
	return Actor{World: w, Actions: actions}
}

func (a *Actor) Act(args ...interface{}) {
	if len(a.Actions) == 0 {
		panic("No actions supplied to actor")
	}
	for _, action := range a.Actions {
		action(a, args...)
	}
}

func (a *Actor) setCoords(xi, yi, zi interface{}) {
	assertAndSet := func(i interface{}) int {
		v, ok := i.(int)
		if !ok {
			panic(fmt.Sprintf("unspecified type for x value %v", xi))
		}
		return v
	}
	a.X = assertAndSet(xi)
	a.Y = assertAndSet(yi)
	a.Z = assertAndSet(zi)
}
