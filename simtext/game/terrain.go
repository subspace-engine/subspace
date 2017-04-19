package game

import (
	"bytes"
	"fmt"
)

type Terrain struct {
	size int
	voxels [5][5][5]TerrainType
	cursor *Cursor
}

type TerrainType uint8

type Position struct {
	x int
	y int
	z int
}

type Cursor Position

const (
	SPACE TerrainType = iota
	GAS
	SAND
	STONE
	ORE
)

func (t *Terrain) GetNameOfTerrainAt(p Position) (terrainName string, err error) {
	terrainType, err := t.GetTerrainTypeAt(p)

	err = nil

	switch terrainType {
	case SPACE:
		terrainName = "space"
	case GAS:
		terrainName = "gas"
	case SAND:
		terrainName = "sand"
	case STONE:
		terrainName = "stone"
	case ORE:
		terrainName = "ore"
	default:
		terrainName = "Unknown"
		// TODO err
	}
	return
}

func (t *Terrain) GetSymbolOfTerrainAt(p Position) (terrainName string, err error) {
	terrainType, err := t.GetTerrainTypeAt(p)

	err = nil

	switch terrainType {
	case SPACE:
		terrainName = "_"
	case GAS:
		terrainName = "."
	case SAND:
		terrainName = "="
	case STONE:
		terrainName = "+"
	case ORE:
		terrainName = "o"
	default:
		terrainName = "?"
		// TODO err = error("Unknown terrain type")
	}
	return
}

func (t *Terrain) DrawnTerrainAtZ(z int) (drawnTerrain string, err error){
	size := t.size
	err = nil
	var byteTerrain bytes.Buffer

	fmt.Println("DrawnTerrainAt Z: ", z)

	for x := 0; x < size; x++ {
		for y:=0; y < size; y++ {
			symbol, _ := t.GetSymbolOfTerrainAt(Position{x,y,z})
			byteTerrain.WriteString(symbol)
		}
		byteTerrain.WriteString("\n")
	}

	drawnTerrain = byteTerrain.String()
	return
}

func (t *Terrain) GetTerrainTypeAt(p Position) (terrainType TerrainType, err error) {
	terrainType = t.voxels[p.x][p.y][p.z]
	err = nil
	return
}

func (w *World) GenerateTerrain() (err error) {
	const size = 5
	mid := size/2
	voxels := [size][size][size]TerrainType{}

	for z := 0; z < mid ; z++ {
		for x := 0; x < size; x++ {
			for y:=0; y < size; y++ {
				voxels[x][y][z] = SAND
			}
		}
	}

	for z := mid; z < size-1 ; z++ {
		for x := 0; x < size; x++ {
			for y:=0; y < size; y++ {
				voxels[x][y][z] = GAS
			}
		}
	}

	for z := size-1; z < size ; z++ {
		for x := 0; x < size; x++ {
			for y:=0; y < size; y++ {
				voxels[x][y][z] = SPACE
			}
		}
	}

	voxels[mid+1][mid][mid] = STONE
	voxels[mid+1][mid+1][mid-1] = STONE


	terrain := &Terrain{size:size, voxels:voxels}
	w.Terrain = terrain
	return nil
}