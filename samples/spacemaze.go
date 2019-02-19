package main

import (
	"fmt"
	"github.com/subspace-engine/subspace/cmd"
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
var lastObj model.Thing

func makeWorld(space world.Space) {
	file, err := os.Open("maze.txt")
	if err != nil {
		panic("Unable to open data file!")
	}
	defer file.Close()
	wall := model.MakePassableThing("wall", "A metal wall", false)
	floor := model.MakePassableThing("floor", "Just the floor", true)
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
	}
}

func playFloorSounds(pos util.Vec3, dir float64, space world.Space) {
	i := snd.PlaySound("step.ogg")
	snd.SetPosition(i, pos)
	pos = pos.Div(space.TileSize())
	left := pos.Add(util.VecFromDirection(dir - math.Pi/2))
	right := pos.Add(util.VecFromDirection(dir + math.Pi/2))
	if space.Encloses(left) {
		time.Sleep(time.Millisecond * 50)
		tile := space.TileAt(int(left.X), int(left.Y), int(left.Z))
		if tile.Passable() {
			i := snd.PlaySound("step.ogg")
			snd.SetPosition(i, left)
		}
	}
	if space.Encloses(right) {
		time.Sleep(time.Millisecond * 20)
		tile := space.TileAt(int(right.X), int(right.Y), int(right.Z))
		if tile.Passable() {
			i := snd.PlaySound("step.ogg")
			snd.SetPosition(i, right)
		}
	}
}

func runTiles() {
	snd.Init()
	cn = con.MakeTextConsole()
	km := cn.Map()
	rand.Seed(time.Now().Unix())
	tiles := world.MakeBasicSpace(100, 100, 100, 1, 20, world.MakeBasicTile(model.MakePassableThing("wall", "just a wall", false)))
	lastObj = nil
	me := model.MakePlayer("you", "As good looking as ever.")
	snd.SetListenerDirection(0)
	me.RegisterAction("describe", func(action model.Action) bool {
		if lastObj != nil {
			cn.Println(lastObj.Description())
		}
		return true
	})
	me.RegisterPrintFunc(cn.Println)
	me.SetPosition(util.Vec3{10, 0, 10})
	me.SetStepSize(0.6)
	snd.SetListenerPosition(me.Position().Add(util.Vec3{0, 1, 0}))
	chair := model.MakeBasicThing("a chair", "Just a basic chair")
	chair.SetPosition(util.Vec3{8, 0, 4})
	me.RegisterAction("bump", func(action model.Action) bool {
		if action.Dobj != nil {
			action.Source.Say("You bumped into " + action.Dobj.Name() + ".")
			lastObj = action.Dobj
			sound := snd.PlaySound("wall.ogg")
			snd.SetPosition(sound, action.Dobj.Position())
		}
		return true
	})
	me.RegisterPostaction("encounter", func(action model.Action) bool {
		if action.Dobj != nil {
			cn.Println("you encountered " + action.Dobj.Name())
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
		snd.SetListenerPosition(me.Position().Add(util.Vec3{0, 1, 0}))
		playFloorSounds(me.Position(), me.Direction(), tiles)
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
	me.RegisterAction("name", func(action model.Action) bool {
		cn.Println(me.Name())
		return true
	})
	parser := cmd.MakeCommandParser(cn, me)
	parser.AddCommand('n', "name", "name")
	parser.AddCommand('i', "forward", "forward")
	parser.MakeKeyAbsolute('i')
	parser.AddCommand(km.KeyUp, "forward", "forward")
	parser.AddCommand(km.KeyRight, "right", "turn right")
	parser.AddCommand(km.KeyLeft, "left", "turn left")
	parser.AddCommand('l', "right", "turn right")
	parser.MakeKeyAbsolute('l')
	parser.AddCommand('j', "left", "turn left")
	parser.MakeKeyAbsolute('j')
	parser.AddCommand('J', "sideleft", "sidestep left")
	parser.MakeKeyAbsolute('J')
	parser.AddCommand('L', "sideright", "sidestep right")
	parser.MakeKeyAbsolute('L')
	parser.AddCommand(km.KeyDown, "reverse", "reverse")
	parser.AddCommand('d', "describe", "describe")
	parser.AddCommand('k', "reverse", "reverse")
	parser.MakeKeyAbsolute('k')
	me.RegisterPostaction("turn left", func(model.Action) bool {
		snd.SetListenerDirection(me.Direction())
		return true
	})
	me.RegisterPostaction("turn right", func(model.Action) bool {
		snd.SetListenerDirection(me.Direction())
		return true
	})
	parser.RunParser()

	cn.Destroy()
	snd.Terminate()
}

func main() {
	runTiles()
}
