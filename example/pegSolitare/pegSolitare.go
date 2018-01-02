package main

import (
	"bufio"
	"fmt"
	"github.com/subspace-engine/subspace/engine/world"
	"os"
)

func main() {
	board := world.NewWorld(9, 9, 1, true, 0)
	reset(board)
	fmt.Println(boardString(board))
	act := board.BuildActor(solitareMove, nil, nil, world.WithTerrain)
	act.X = 4
	act.Y = 4
	act.CoordsSet = true
	act2 := board.BuildActor(pegsRemaining, nil, nil, world.WithTerrain)
	act2.CoordsSet = true

	// startServer()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "u":
			act.Y--
		case "r":
			act.X++
		case "d":
			act.Y++
		case "l":
			act.X--
		case "c":
			fmt.Printf("%s, %d, %d\n", squareString(act.World.Terrain[act.X][act.Y][0]), act.X, act.Y)
		case "B":
			fmt.Println(boardString(board))
		case "U":
			act.Act(dirUp)
		case "R":
			act.Act(dirRight)
		case "D":
			act.Act(dirDown)
		case "L":
			act.Act(dirLeft)
		case "t":
			act2.Act()
		}
	}
}

const (
	dirUp    uint8 = 1
	dirRight uint8 = 2
	dirDown  uint8 = 3
	dirLeft  uint8 = 4
)

func solitareMove(act *world.Actor, args ...interface{}) {
	if len(args) < 1 {
		panic(fmt.Sprintf("1 arguments is required, got %d", len(args)))
	}
	direction, ok := args[0].(uint8)
	if !ok {
		panic("could not convert argument to direction")
	}
	if act.World.Terrain[act.X][act.Y][0] != peg {
		fmt.Println("square is not a peg")
		return
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
		if act.X > 5 {
			fmt.Println("illegal move")
			return
		}
		b = &act.World.Terrain[act.X+1][act.Y][0]
		c = &act.World.Terrain[act.X+2][act.Y][0]
	case dirDown:
		if act.Y > 5 {
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
	if *b != peg {
		fmt.Println("No peg to jump over")
		return
	}
	if *c != blank {
		fmt.Println("No empty square to land on")
		return
	}
	act.World.Terrain[act.X][act.Y][0] = blank
	*b = blank
	*c = peg
	act.World.GlobalObjects["count"] = act.World.GlobalObjects["count"].(int) - 1
}

func pegsRemaining(act *world.Actor, args ...interface{}) {
	v, ok := act.World.GlobalObjects["count"]
	if !ok {
		panic("global variable count not found")
	}
	i := v.(int)
	fmt.Printf("%d pegs\n", i)
}
