package main

import (
	"fmt"
	"github.com/subspace-engine/subspace/con"
	"github.com/subspace-engine/subspace/world"
	"github.com/subspace-engine/subspace/world/model"
	"math/rand"
	"time"
)

var cn con.Console

func makeWorld(space world.Space) {
	wall := model.MakePassableThing("wall", "A metal wall", false)
	floor := model.MakePassableThing("floor", "Just the floor", true)
	for i := 0; i < 5000; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 5000; k++ {
				if rand.Float64() < 0.001 {
					space.SetTile(i, j, k, world.MakeBasicTile(wall))
				} else {
					space.SetTile(i, j, k, world.MakeBasicTile(floor))
				}
			}
		}
	}
}

func runTiles() {
	cn = con.MakeTextConsole()
	km := cn.Map()
	proc := cn.MakeEventProc()
	rand.Seed(time.Now().Unix())
	tiles := world.MakeDefaultSpace(5000, 4, 5000)
	me := model.MakeBasicThing("you", "As good looking as ever.")
	me.SetX(1)
	me.SetY(0)
	me.SetZ(1)
	me.RegisterReaction("bump", func(action model.Action) int {
		if action.Dobj != nil {
			cn.Println("You bumped into " + action.Dobj.Name() + ".")
		}
		return 0
	})
	me.RegisterReaction("step", func(action model.Action) int {
		if action.Dobj != nil {
			cn.Println(action.Dobj.Name())
		}
		return 0
	})
	makeWorld(tiles)
	tiles.Add(1, 0, 1, me)
	running := true
	proc.SetKeyDown(func(key int) {
		switch key {
		case 27:
			running = false
		case ' ':
			cn.Println(fmt.Sprintf("%.1f, %.1f, %.1f\n", me.X(), me.Y(), me.Z()))
		case km.KeyUp:
			tiles.Move(me, 0, 0, -1)
		case km.KeyDown:
			tiles.Move(me, 0, 0, 1)
		case km.KeyLeft:
			tiles.Move(me, -1, 0, 0)
		case km.KeyRight:
			tiles.Move(me, 1, 0, 0)
		}
	})
	for running {
		time.Sleep(3000000)
		proc.Pump()
	}

	cn.Destroy()
}

func main() {
	runTiles()
}
