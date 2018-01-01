package main

import (
	"fmt"
	"github.com/subspace-engine/subspace/engine/world"
)

func main() {
	board := world.NewWorld(9, 9, 1, true, 0)
	reset(board)
	fmt.Println(boardString(board))
	act := board.BuildActor(solitareMove, nil, nil, world.WithTerrain)
	act2 := board.BuildActor(pegsRemaining, nil, nil, world.WithTerrain)
	act2.CoordsSet = true

	// startServer()
	act.Act()
}

const (
	dirUp    = 1
	dirRight = 2
	dirDown  = 3
	dirLeft  = 4
)

func solitareMove(act *world.Actor, args ...interface{}) {
	if len(args) < 1 {
		panic(fmt.Sprintf("3 arguments are required, got %d", len(args)))
	}
	direction, ok := args[0].(uint8)
	if !ok {
		panic("could not convert argument to direction")
	}
	var b, c *world.TerrainType
	switch direction {
	case dirUp:
		if act.Y < 2 {
			fmt.Println("illegal move")
			return
		}
		b = &act.World.Terrain[act.X][act.Y-1][0]
		c = &act.World.Terrain[act.X][act.Y-2][0]
	case dirRight:
		if int(act.World.YSize)-act.Y < 3 {
			fmt.Println("illegal move")
			return
		}
		b = &act.World.Terrain[act.X+1][act.Y][0]
		c = &act.World.Terrain[act.X+2][act.Y][0]
	case dirDown:
		if int(act.World.XSize)-act.X < 3 {
			fmt.Println("illegal move")
			return
		}
		b = &act.World.Terrain[act.X][act.Y+1][0]
		c = &act.World.Terrain[act.X][act.Y+2][0]

	case dirLeft:
		if act.X < 2 {
			fmt.Println("illegal move")
			return
		}
		b = &act.World.Terrain[act.X-1][act.Y][0]
		c = &act.World.Terrain[act.X-2][act.Y][0]
	default:
		panic("unknown direction")
	}
	act.World.Terrain[act.X][act.Y][act.Z] = blank
	*b = blank
	*c = peg
}

func pegsRemaining(act *world.Actor, args ...interface{}) {
	v, ok := act.World.GlobalObjects["count"]
	if !ok {
		panic("global variable count not found")
	}
	i := v.(int)
	fmt.Printf("%d pegs\n", i)
}
