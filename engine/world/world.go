package world

// World, a type intended to hold all the objects of a game and how they relate to each other.

import (
	"fmt"
)

// UnsetDimention indicates when a world does not use a particular dimention, e.g. card games will use this for X, Y and ZSize, while 2-dimentional worlds will only use it for the ZSize
const UnsetDimention = 0.0

type TerrainType uint8

const TerrainUnset = 0

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

func NewTiles(x, y, z, sf int) [][][]*ObjectTile {
	o1 := scaleSize(x, sf)
	o2 := scaleSize(y, sf)
	o3 := scaleSize(z, sf)
	tiles := make([][][]*ObjectTile, o1, o1)
	for i := range tiles {
		tiles[i] = make([][]*ObjectTile, o2, o2)
		for j := 0; j < o3; j++ {
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

type World struct {
	XSize float64
	YSize float64
	ZSize float64
	SF    int // scale factor of sectors to object tiles

	//below types can be optimised later to not be multidimentional arrays
	// Terrain represents any cube of world space that requires unique feature types, setting this to nil disables associated logic
	Terrain       [][][]TerrainType
	ObjectTiles   [][][]*ObjectTile
	GlobalObjects map[string]interface{}
}

func NewWorld(x, y, z int, hasTerrain bool, sf int) *World {
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
	return &World{XSize: float64(x), YSize: float64(y), ZSize: float64(z), SF: sf, Terrain: terrain, ObjectTiles: tiles, GlobalObjects: make(map[string]interface{})}
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

func (w *World) AddObject(x, y, z int, obj interface{}) error {
	//todo, bounds check
	const offset = 0 // assume scale 1 for now, todo other sfs
	w.ObjectTiles[x][y][z].Add(Tile{offset, offset, offset, obj})
	return nil
}

type ActorFunc func(actor *Actor, args ...interface{})

type ActWith int

const (
	withNil     ActWith = 0
	WithTerrain ActWith = 1
	WithTile    ActWith = 2
	WithObjects ActWith = 3
	WithGlobals ActWith = 4
)

type Actor struct {
	World        *World
	TerrainActor ActorFunc
	TileActor    ActorFunc
	ObjectActor  ActorFunc
	GlobalActor  ActorFunc
	Order        []ActWith
	X            int
	Y            int
	Z            int
	Objects      []*Tile
	Terrain      TerrainType
	CoordsSet    bool
}

// pass actors by value to reuse them without changing the basis Actor
// I think this is what we usually want
// pass using a pointer to mutate
func (w *World) BuildActor(terrainActor, tileActor, objectActor ActorFunc, order ...ActWith) Actor {
	return Actor{World: w,
		TerrainActor: terrainActor,
		TileActor:    tileActor,
		ObjectActor:  objectActor,
		Order:        order}
}

func (a *Actor) Act(args ...interface{}) {
	if len(a.Order) == 0 {
		panic("Action order must be specified")
	}
	setAttrs := func() {
		if !a.CoordsSet {
			// first 3 arguments must be an x, y, z coord
			if len(args) >= 3 {
				a.setCoords(args[0], args[1], args[2])
				args = args[3:]
				if a.World.Terrain != nil {
					tt := a.World.Terrain[a.X][a.Y][a.Z]
					if tt == TerrainUnset {
						panic(fmt.Sprintf("found unspecified terrain at %d, %d, %d", a.X, a.Y, a.Z))
					}
					a.Terrain = tt
				}
				return
			}
			panic(fmt.Sprintf("only %d arguments found: %v", len(args), args))
		}
	}

	for _, fType := range a.Order {
		switch fType {
		case WithTerrain:
			setAttrs()
			a.TerrainActor(a, args)
		case WithTile:
			setAttrs()
			a.TileActor(a, args)
		case WithObjects:
		default:
			panic(fmt.Sprintf("Unspecified action order %d\n", fType))
		}
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
