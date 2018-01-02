package main

import (
	"fmt"
	"github.com/subspace-engine/subspace/engine/world"
)

type square world.TerrainType

const (
	unset = 0
	edge  = 1
	blank = 2
	peg   = 3
)

func reset(w *world.World) {
	w.GlobalObjects["count"] = 32
	const middle = 4
	for i := 0; i < 9; i++ {
		st := typeAt(i) // not x y index based

		x := i % 3
		y := i / 3
		// top left
		w.Terrain[x+1][y+1][0] = st

		//top right
		w.Terrain[7-x][y+1][0] = st

		// bottom left
		w.Terrain[x+1][7-y][0] = st

		// bottom right
		w.Terrain[7-x][7-y][0] = st

		// middle pegs
		if !(i >= 7 || i == 3) {
			w.Terrain[i+1][middle][0] = peg
			w.Terrain[middle][i+1][0] = peg
		}

		//edge
		w.Terrain[i][0][0] = edge
		w.Terrain[i][8][0] = edge
		w.Terrain[8][i][0] = edge
		w.Terrain[0][i][0] = edge

	}
	w.Terrain[middle][middle][0] = blank
}

func typeAt(index int) world.TerrainType {
	switch index {
	case 0, 1, 3, 4:
		return edge
	case 2, 5, 6, 7, 8:
		return peg
	default:
		fmt.Printf("index %d is unset\n", index)
		return unset
	}
}

func boardString(w *world.World) string {
	s := ""
	for y := 0; y < 9; y++ {
		s = s + "\n"
		for x := 0; x < 9; x++ {
			t := w.Terrain[x][y][0]
			sq := squareString(t)
			if sq == "" {
				panic(fmt.Sprintf("found tile of type %d at index %d, %d\n", t, x, y))
			}
			s = s + sq
		}
	}
	return s
}

func squareString(t world.TerrainType) string {
	switch t {
	case edge:
		return "E"
	case blank:
		return "B"
	case peg:
		return "P"
	case 0:
		return "0"
	}
	return ""
}
