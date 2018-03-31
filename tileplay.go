package main

import (
	"github.com/nsf/termbox-go"
	"github.com/subspace-engine/subspace/world"
	"github.com/subspace-engine/subspace/world/model"
	"github.com/subspace-engine/subspace/term"
	"fmt"
)

func makeWorld(space world.Space) {
	for i:=40; i <61; i++ {
j:=0
for k:=40; k<61; k++ {
				space.SetTile(i,j,k,world.MakeBasicTile(world.Floor))
}
		}

}

func runTiles() {
	term.Init()
	tiles := world.MakeDefaultSpace(100,100,100)
	me := model.MakeThing("you", "As good looking as ever.")
me.SetX(50)
	me.SetY(0)
	me.SetZ(50)
	makeWorld(tiles)
	tiles.Add(50,0,50,&me)
	loop:
	for {
		switch term.Read() {
				case 27:
				break loop
			case ' ':
				term.Println(fmt.Sprintf("%f, %f, %f\n",me.X(), me.Y(), me.Z()))
				case term.KeyUp():
				tiles.Move(&me, 0, 0, -1)
				term.Println(tiles.GetTile(&me))
				case term.KeyDown():
				tiles.Move(&me,0,0,1)
				term.Println(tiles.GetTile(&me))
				case term.KeyLeft():
				tiles.Move(&me,-1,0,0)
				term.Println(tiles.GetTile(&me))
				case term.KeyRight():
				tiles.Move(&me,1,0,0)
				term.Println(tiles.GetTile(&me))
			}
		}




term.Terminate()
}

func main() {
	runTiles()
}
