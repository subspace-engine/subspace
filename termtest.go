package main

import "github.com/subspace-engine/subspace/term"

func main() {
term.Init()
term.Println("Hi, welcome to this example! Press q to quit.")
for {
c := term.Read()
if c == 'q' {
break
}
if c == 58 {
term.Print(":")
text := term.Readln()
term.Println("You entered " + text)
}
}
term.Terminate()
}
