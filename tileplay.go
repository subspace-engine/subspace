package main

import (
"github.com/subspace-engine/subspace/con"
	"github.com/subspace-engine/subspace/world"
	"github.com/subspace-engine/subspace/world/model"
	"fmt"
)

func makeWorld(space world.Space) {
	for i:=40; i <61; i++ {
j:=0
for k:=40; k<61; k++ {
	space.SetTile(i,j,k,world.MakeBasicTile(model.MakePassableThing("floor", "Just the floor", true)))
}
		}

}

func runTiles() {
	con := con.MakeTextConsole()
	km :=con.Map()
	proc :=con.MakeEventProc()
	tiles := world.MakeDefaultSpace(100,100,100)
	me := model.MakeBasicThing("you", "As good looking as ever.")
me.SetX(50)
	me.SetY(0)
	me.SetZ(50)
	makeWorld(tiles)
	tiles.Add(50,0,50,me)
running :=true
		proc.SetKeyDown(func (key int) {
		switch (key) {
				case 27:
				running=false
			case ' ':
				con.Println(fmt.Sprintf("%f, %f, %f\n",me.X(), me.Y(), me.Z()))
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
	proc.Pump()
}


con.Destroy()
		}


func main() {
	runTiles()
}
