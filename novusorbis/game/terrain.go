package game

type Terrain struct {
	size int
	voxels [][][]TerrainType
	TerrainToString map[TerrainType]string
	TerrainToSymbol map[TerrainType]string
	TerrainToOpacity map[TerrainType]TerrainOpacity
}

type TerrainType uint8
type TerrainOpacity bool

const (
	SPACE TerrainType = iota
	GAS
	SAND
	STONE
	ORE
	UNKNOWN
)

const (
	OPAQUE TerrainOpacity = true
	CLEAR = false
)

func (t *Terrain) SetUpTerrainConversions() (err error) {
	terrainToString := map[TerrainType]string{SPACE : "space",GAS : "gas",SAND : "sand",
		STONE : "stone",ORE : "ore",UNKNOWN : "unknown"}
	t.TerrainToString = terrainToString

	terrainToSymbol := map[TerrainType]string{SPACE: ".",GAS: "-",SAND: "~",
		STONE: "S",ORE: "o",UNKNOWN: "_"}
	t.TerrainToSymbol = terrainToSymbol

	terrainToOpacity := map[TerrainType]TerrainOpacity{SPACE: CLEAR, GAS: CLEAR, SAND: OPAQUE,
		STONE: OPAQUE, ORE: OPAQUE, UNKNOWN: OPAQUE}
	t.TerrainToOpacity = terrainToOpacity

	err = nil
	return
}

func (t *Terrain) GetNameOfTerrainAt(p *Position) (terrainName string, err error) {
	terrainType, err := t.GetTerrainTypeAt(p)
	terrainName = t.TerrainToString[terrainType]
	return
}

func (t *Terrain) GetSymbolOfTerrainAt(p *Position) (terrainChar string, err error) {
	terrainType, err := t.GetTerrainTypeAt(p)
	terrainChar = t.TerrainToSymbol[terrainType]
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
	terrain.SetUpTerrainConversions()
	w.Terrain = terrain
	return nil
}