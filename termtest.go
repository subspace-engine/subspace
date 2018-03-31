package main

import "github.com/subspace-engine/subspace/term"
import "fmt"

func main() {
term.Init()
term.Println("Hi, welcome to this example! Press q to quit.")
var x int=0
var y int = 0
running :=true
for running {
switch term.Read() {
case 'q':
running=false
case ':':
term.Print(":")
text := term.Readln()
term.Println("You entered " + text)
case term.KeyUp():
y+=1
term.Println(fmt.Sprintf("X:%d, y:%d", x, y))
case term.KeyDown():
y-=1
term.Println(fmt.Sprintf("X:%d, y:%d", x, y))
case term.KeyLeft():
x-=1
term.Println(fmt.Sprintf("X:%d, y:%d", x, y))
case term.KeyRight():
x+=1
term.Println(fmt.Sprintf("X:%d, y:%d", x, y))
}
}
term.Terminate()
}
