package main

import (
"github.com/subspace-engine/subspace/con"
	"github.com/subspace-engine/subspace/world"
	"github.com/subspace-engine/subspace/world/model"
	"os"
	"time"
	"fmt"
)

func makeWorld(space world.Space) {
	file, err :=os.Open("maze.txt")
	if err!=nil {
		panic("Unable to open data file!")
	}
	defer file.Close()
	wall :=model.MakePassableThing("wall", "A metal wall", false)
	floor := model.MakePassableThing("floor", "Just the floor", true)
	for i:=0; i < 100; i++ {
		for j:=0; j < 100; j++ {
			for k :=0; k < 100; k++ {
				space.SetTile(i,j,k,world.MakeBasicTile(wall))
			}
		}
	}
reader:
		for {
		x,y,w,h :=0,0,0,0
		_, err := fmt.Fscanf(file, "%d, %d, %d, %d\n", &x, &y,&w,&h)
			if err !=nil {
//				fmt.Println("unable to read coords")
//				fmt.Println(err)
			break reader
			}
			fmt.Println("Got coords")
	for k:=y; k <y+h; k++ {
j:=0
		for i:=x; i<x+w; i++ {
			symbol :=' '
			_, err := fmt.Fscanf(file, "%c", &symbol)
			if err!=nil {
				fmt.Println("Unable to read letter")
				break reader
			}
			switch (symbol) {
				case 'x':
				space.SetTile(i,j,k,world.MakeBasicTile(floor))
				fmt.Printf("Painting floor at position %d, %d\n", i, k)
			case 'w':
								space.SetTile(i,j,k,world.MakeBasicTile(wall))
			}

		}
					_, err = fmt.Fscanf(file, "\n")
			if err!=nil {
				break reader
			}
		}
	}
}

func runTiles() {
	con := con.MakeTextConsole()
	km :=con.Map()
	proc :=con.MakeEventProc()
	tiles := world.MakeDefaultSpace(100,100,100)
	me := model.MakeBasicThing("you", "As good looking as ever.")
me.SetX(1)
	me.SetY(0)
	me.SetZ(1)
	me.RegisterAction("bump", func(action model.Action) int {
		if action.Dobj!=nil {
			con.Println("You bumped into " + action.Dobj.Name()+".")
		}
		return 0
	})
	makeWorld(tiles)
	tiles.Add(1,0,1,me)
running :=true
		proc.SetKeyDown(func (key int) {
		switch (key) {
				case 27:
				running=false
			case ' ':
				con.Println(fmt.Sprintf("%.1f, %.1f, %.1f\n",me.X(), me.Y(), me.Z()))
				case km.KeyUp:
				tiles.Move(me, 0, 0, -1)
			con.Println(tiles.GetTile(me).String())
				case km.KeyDown:
				tiles.Move(me,0,0,1)
			con.Println(tiles.GetTile(me).String())
				case km.KeyLeft:
				tiles.Move(me,-1,0,0)
			con.Println(tiles.GetTile(me).String())
				case km.KeyRight:
				tiles.Move(me,1,0,0)
			con.Println(tiles.GetTile(me).String())
			}
		})
	for running {
		time.Sleep(3000000)
	proc.Pump()
}


con.Destroy()
		}


func main() {
	runTiles()
}
