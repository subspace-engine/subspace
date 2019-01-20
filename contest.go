package main

import(
	"fmt"
	"github.com/subspace-engine/subspace/con"
)


func main () {
	con := con.MakeTextConsole()
loop :=true
	km :=con.Map()
	con.Println("Hello, world!")
	proc :=con.MakeEventProc()
	proc.SetKeyDown(func (key int) {
		switch(key) {
		case 'q':
			con.Println("Q!")
			loop=false
		case km.KeyUp:
			con.Println("UP!")
		default:
			con.Println(fmt.Sprintf("%d", key))
		}
	})
	for loop {
		proc.Pump()
	}
}
