package main

import (
	"fmt"
	"github.com/subspace-engine/subspace/con"
	"github.com/subspace-engine/subspace/util"
	"github.com/subspace-engine/subspace/world"
	"github.com/subspace-engine/subspace/world/model"
	"math"
	"math/rand"
	"os"
	"time"
)

var cn con.Console

func isOpen(tile world.Tile) bool {
	return tile.Name() == "floor" || tile.Name() == "intersection" || tile.Name() == "Tee intersection" || tile.Name() == "doorway"
}

func printDirection(direction float64) {
	cn.Println(fmt.Sprintf("%d degrees", int(direction/(math.Pi*2)*360)))
}

func makeShip(ship world.Space, tiles world.Space) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 10; k++ {
				ship.SetTile(i, j, k, world.MakeBasicTile(model.MakePassableThing("ship wall", "a metal wall", false)))
			}
		}
	}
	for i := 3; i < 9; i++ {
		ship.SetTile(i, 0, 2, world.MakeBasicTile(model.MakePassableThing("floor", "Metal floor", true)))
		ship.SetTile(4, 0, i, world.MakeBasicTile(model.MakePassableThing("floor", "Metal floor", true)))
	}
	exit := model.MakePassableThing("Exit", "The exit", true)
	exit.RegisterAction("step", func(action model.Action) int {
		if action.Source != nil {
			ship.Remove(action.Source)
			action.Source.SetPosition(util.Vec3{1, 0, 2})
			tiles.Add(action.Source)
		}
		return 0
	})
	ship.SetTile(4, 0, 1, world.MakeBasicTile(exit))
}

func makeWorld(space world.Space) {
	file, err := os.Open("maze.txt")
	if err != nil {
		panic("Unable to open data file!")
	}
	defer file.Close()
	wall := model.MakePassableThing("wall", "A metal wall", false)
	floor := model.MakePassableThing("floor", "Just the floor", true)
	intersection := model.MakePassableThing("intersection", "Intersaction between two paths", true)
	tee := model.MakePassableThing("tee intersection", "Tee intersection with a branching path", true)
	doorway := model.MakePassableThing("doorway", "A doorway into another room", true)
	hole := model.MakePassableThing("Hole", "A hole in dhe ground", true)
	hole.RegisterAction("step", func(action model.Action) int {
		cn.Println("You stepped into a hole. You die.")
		return 0
	})
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			for k := 0; k < 100; k++ {
				space.SetTile(i, j, k, world.MakeBasicTile(wall))
			}
		}
	}
reader:
	for {
		x, y, w, h := 0, 0, 0, 0
		_, err := fmt.Fscanf(file, "%d, %d, %d, %d\n", &x, &y, &w, &h)
		if err != nil {
			//				fmt.Println("unable to read coords")
			//				fmt.Println(err)
			break reader
		}
		for k := y; k < y+h; k++ {
			j := 0
			for i := x; i < x+w; i++ {
				symbol := ' '
				_, err := fmt.Fscanf(file, "%c", &symbol)
				if err != nil {
					fmt.Println("Unable to read letter")
					break reader
				}
				switch symbol {
				case 'x':
					space.SetTile(i, j, k, world.MakeBasicTile(floor))
				case 'w':
					space.SetTile(i, j, k, world.MakeBasicTile(wall))
				}

			}
			_, err = fmt.Fscanf(file, "\n")
			if err != nil {
				break reader
			}
		}
		for k := y + 1; k < y+h-1; k++ {
			j := 0
			for i := x + 1; i < x+w-1; i++ {
				if isOpen(space.TileAt(i, j, k)) {
					horizontal := 0
					vertical := 0
					corners := 0
					for m := i - 1; m <= i+1; m++ {
						if isOpen(space.TileAt(m, j, k)) {
							horizontal++
						}
					}
					for m := k - 1; m <= k+1; m++ {
						if isOpen(space.TileAt(i, j, m)) {
							vertical++
						}
					}
					for m := 0; m < 4; m++ {
						xmin := (m%2)*2 - 1
						ymin := (m/2)*2 - 1

						if isOpen(space.TileAt(i+xmin, j, k+ymin)) {
							corners++
						}
					}
					if horizontal == 3 && vertical == 3 && corners == 0 {
						space.SetTile(i, j, k, world.MakeBasicTile(intersection))
					} else if (horizontal+vertical) == 5 && corners <= 1 {
						space.SetTile(i, j, k, world.MakeBasicTile(tee))
					} else if (horizontal == 3 || vertical == 3) && corners == 4 {
						space.SetTile(i, j, k, world.MakeBasicTile(doorway))
					}

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
	tiles := world.MakeBasicSpace(100, 100, 100, 1, 20, world.MakeBasicTile(model.MakePassableThing("wall", "just a wall", false)))
	ship := world.MakeDefaultSpace(10, 10, 10)
	ship.SetName("A ship")
	ship.SetDescription("Looks like a space ship of some kind")
	ship.SetPassable(true)
	makeShip(ship, tiles)
	ship.RegisterAction("step", func(action model.Action) int {
		cn.Println("Entering ship")
		if action.Source != nil {
			cn.Println(fmt.Sprintf("%s are being transported", action.Source.Name()))
			tiles.Remove(action.Source)
			action.Source.SetPosition(util.Vec3{5, 0, 2})
			ship.Add(action.Source)
		}
		return 0
	})
	me := model.MakeMobileThing("you", "As good looking as ever.")
	me.SetPosition(util.Vec3{10, 0, 10})
	chair := model.MakeBasicThing("a chair", "Just a basic chair")
	chair.SetPosition(util.Vec3{8, 0, 4})
	me.RegisterAction("bump", func(action model.Action) int {
		if action.Dobj != nil {
			cn.Println("You bumped into " + action.Dobj.Name() + ".")
		}
		return 0
	})
	me.RegisterPostaction("encounter", func(action model.Action) int {
		if action.Dobj != nil {
			cn.Println(fmt.Sprintf("You encountered %s", action.Dobj.Name()))
		}
		return 0
	})
	me.RegisterAction("step", func(action model.Action) int {
		if action.Dobj != nil {
			cn.Println(action.Dobj.Name())
		}
		return 0
	})
	makeWorld(tiles)
	tiles.SetTile(2, 0, 2, ship)
	tiles.Add(me)
	tiles.Add(chair)
	running := true
	proc.SetKeyDown(func(key int) {
		switch key {
		case 27:
			running = false
		case ' ':
			cn.Println(fmt.Sprintf("%.1f, %.1f, %.1f\n", me.Position().X, me.Position().Y, me.Position().Z))
		case km.KeyUp:
			me.Act(model.Action{me, "move", nil, nil})
		case km.KeyDown:
			tiles.Move(me, util.Vec3{0, 0, 1})
		case km.KeyLeft:
			me.SetDirection(me.Direction() - math.Pi/2)
			me.SetDirection(math.Mod(math.Pi*2+me.Direction(), (math.Pi * 2)))
			printDirection(me.Direction())
		case km.KeyRight:
			me.SetDirection(me.Direction() + math.Pi/2)
			me.SetDirection(math.Mod(math.Pi*2+me.Direction(), (math.Pi * 2)))
			printDirection(me.Direction())
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
