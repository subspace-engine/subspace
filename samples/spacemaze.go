package main

import (
	"fmt"
	"github.com/subspace-engine/subspace/con"
	"github.com/subspace-engine/subspace/snd"
	"github.com/subspace-engine/subspace/util"
	"github.com/subspace-engine/subspace/world"
	"github.com/subspace-engine/subspace/world/model"
	"math"
	"math/rand"
	"os"
	"time"
)

var cn con.Console

func getTileAhead(tiles world.Space, thing model.MobileThing) world.Tile {
	pos := thing.Position()
	pos.X += math.Sin(thing.Direction())
	pos.Z -= math.Cos(thing.Direction() + 0.000001)
	tile := tiles.TileAt(int(pos.X), int(pos.Y), int(pos.Z))
	return tile
}

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
	exit.RegisterAction("step", func(action model.Action) bool {
		if action.Source != nil {
			ship.Remove(action.Source)
			action.Source.SetPosition(util.Vec3{1, 0, 2})
			tiles.Add(action.Source)
			action.Source.Say(fmt.Sprintf("Tile at 2,0,1: %t\n", tiles.TileAt(2, 0, 1)))
		}
		return true
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
	hole.RegisterAction("step", func(action model.Action) bool {
		cn.Println("You stepped into a hole. You die.")
		return true
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

func playFloorSounds(pos util.Vec3, dir float64, space world.Space) {
	i := snd.PlaySound("footstep.ogg")
	snd.SetPosition(i, pos)
	left := pos.Add(util.VecFromDirection(dir - math.Pi/2))
	right := pos.Add(util.VecFromDirection(dir + math.Pi/2))
	time.Sleep(time.Millisecond * 100)
	if space.Encloses(left) {
		tile := space.TileAt(int(left.X), int(left.Y), int(left.Z))
		if tile.Passable() {
			i := snd.PlaySound("footstep.wav")
			snd.SetPosition(i, left)
		}
	}
	if space.Encloses(right) {
		tile := space.TileAt(int(right.X), int(right.Y), int(right.Z))
		if tile.Passable() {
			i := snd.PlaySound("footstep.wav")
			snd.SetPosition(i, right)
		}
	}
}

func runTiles() {
	snd.Init()
	//go snd.PlaySound("/home/rkruger/audio/dingdong.wav")
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
	ship.RegisterAction("step", func(action model.Action) bool {
		if action.Source == nil {
			return false
		}
		action.Source.Say("Entering ship")
		action.Source.Say(fmt.Sprintf("%s are being transported", action.Source.Name()))
		tiles.Remove(action.Source)
		action.Source.SetPosition(util.Vec3{5, 0, 2})
		ship.Add(action.Source)
		return true
	})
	me := model.MakePlayer("you", "As good looking as ever.")
	snd.SetListenerDirection(0)
	me.RegisterPrintFunc(cn.Println)
	//	me := model.MakeMobileThing("you", "As good looking as ever.")
	me.SetPosition(util.Vec3{10, 0, 10})
	chair := model.MakeBasicThing("a chair", "Just a basic chair")
	chair.SetPosition(util.Vec3{8, 0, 4})
	me.RegisterAction("bump", func(action model.Action) bool {
		if action.Dobj != nil {
			action.Source.Say("You bumped into " + action.Dobj.Name() + ".")
			sound := snd.PlaySound("wall.wav")
			snd.SetPosition(sound, action.Dobj.Position())
		}
		return true
	})
	me.RegisterPostaction("encounter", func(action model.Action) bool {
		if action.Dobj != nil {
			cn.Println("you encountered %s." + action.Dobj.Name())
		}
		return true
	})
	me.RegisterPostaction("move", func(action model.Action) bool {
		if action.Source == nil {
			return true
		}
		loc := action.Source.Location()
		if loc == nil {
			return true
		}
		snd.SetListenerPosition(me.Position())
		go playFloorSounds(me.Position(), me.Direction(), tiles)
		return true
	})
	makeWorld(tiles)
	machine := world.MakeBasicTile(model.MakePassableThing("Machine", "Some kind of machine", false))
	tiles.SetTile(2, 0, 2, machine)
	machinesound := snd.PlaySound("machine.wav")
	snd.SetLooping(machinesound, true)
	snd.SetPosition(machinesound, machine.Position())
	tiles.Add(me)
	tiles.Add(chair)
	running := true
	proc.SetKeyDown(func(key int) {
		switch key {
		case 27:
			running = false
		case 'd':
			me.Say(fmt.Sprintf("Direction %d, radians %fPi\n", int(me.Direction()/(math.Pi*2)*360), me.Direction()/math.Pi))
		case 'l':
			tile := getTileAhead(tiles, me)
			me.Say(fmt.Sprintf("Ahead: %s at %f, %f, %f\n", tile.Name(), tile.Position().X, tile.Position().Y, tile.Position().Z))
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
			snd.SetListenerDirection(me.Direction())
		case km.KeyRight:
			me.SetDirection(me.Direction() + math.Pi/2)
			me.SetDirection(math.Mod(math.Pi*2+me.Direction(), (math.Pi * 2)))
			snd.SetListenerDirection(me.Direction())
			printDirection(me.Direction())
		}
	})
	for running {
		time.Sleep(3000000)
		proc.Pump()
	}

	cn.Destroy()
	snd.Terminate()
}

func main() {
	runTiles()
}
