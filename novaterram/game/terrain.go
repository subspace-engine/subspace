package game

type Terrain struct {
	size int
	voxels [][][]TerrainType
}

type TerrainType uint8


const (
	SPACE TerrainType = iota
	GAS
	SAND
	STONE
	ORE
	UNKNOWN
)

func (t *Terrain) GetNameOfTerrainAt(p *Position) (terrainName string, err error) {
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
	case UNKNOWN:
		terrainName = "unknown"
	default:
		terrainName = "wrong"
		// TODO err
	}
	return
}

func (t *Terrain) GetSymbolOfTerrainAt(p *Position) (terrainChar string, err error) {
	terrainType, err := t.GetTerrainTypeAt(p)

	err = nil

	switch terrainType {
	case SPACE:
		terrainChar = "."
	case GAS:
		terrainChar = "-"
	case SAND:
		terrainChar = "~"
	case STONE:
		terrainChar = "{"
	case ORE:
		terrainChar = "o"
	case UNKNOWN:
		terrainChar = "_"
	default:
		terrainChar = "??"
		// TODO err = error("Unknown terrain type")
	}
	return
}

func (t *Terrain) GetTerrainTypeAt(p *Position) (terrainType TerrainType, err error) {
	terrainType = t.voxels[p.z][p.y][p.x]
	err = nil
	return
}

func (w *World) GenerateTerrain() (err error) {
	var size = w.Size
	mid := size/2
	voxels := make([][][]TerrainType, size)
	for z := 0; z < size; z++ {
		voxels[z] = make([][]TerrainType, size)
		for y := 0; y < size; y++ {
			voxels[z][y] = make([]TerrainType, size)
		}
	}


	for z := 0; z < mid ; z++ {
		for x := 0; x < size; x++ {
			for y:=0; y < size; y++ {
				voxels[z][y][x] = STONE
			}
		}

	}

	for z := mid; z < size-1 ; z++ {
		for x := 0; x < size; x++ {
			for y:=0; y < size; y++ {
				voxels[z][y][x] = GAS
			}
		}
	}

	for z := size-1; z < size ; z++ {
		for x := 0; x < size; x++ {
			for y:=0; y < size; y++ {
				voxels[z][y][x] = SPACE
			}
		}
	}

	voxels[mid][mid][mid+1] = STONE
	voxels[mid+1][mid+1][mid-1] = STONE

	terrain := &Terrain{size:size, voxels:voxels}
	w.Terrain = terrain
	return nil
}