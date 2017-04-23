package main

import (
	"github.com/nsf/termbox-go"
	"github.com/subspace-engine/subspace/world"
	"github.com/subspace-engine/subspace/world/model"
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
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	tiles := world.MakeDefaultSpace(100,100,100)
	me := model.MakeThing("you", "As good looking as ever.")
me.SetX(50)
	me.SetY(0)
	me.SetZ(50)
	makeWorld(tiles)
	tiles.Add(50,0,50,&me)
	loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyEnd:
				fmt.Printf("%f, %f, %f\n",me.X(), me.Y(), me.Z())
			case termbox.KeyArrowUp:
				tiles.Move(&me, 0, 0, -1)
				fmt.Println(tiles.GetTile(&me))
							case termbox.KeyArrowDown:
				tiles.Move(&me,0,0,1)
				fmt.Println(tiles.GetTile(&me))
			case termbox.KeyArrowLeft:
				tiles.Move(&me,-1,0,0)
				fmt.Println(tiles.GetTile(&me))
			case termbox.KeyArrowRight:
				tiles.Move(&me,1,0,0)
				fmt.Println(tiles.GetTile(&me))
			}
		}

		termbox.Flush()
	}

	defer termbox.Close()
}

func main() {
	runTiles()
}
