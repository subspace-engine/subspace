package game_test

import (
	"testing"
	"github.com/subspace-engine/subspace/xenoterra/game"
)

type InputOutput interface {
	Print(s string)
	Println(s string)
	Read() (s string)
}


type TestInputOutput struct{
	Iter int
}

func (t *TestInputOutput) Print(s string) {
	// Do nothing
}

func (t *TestInputOutput) Println(s string) {
	// Do nothing
}

func (t *TestInputOutput) Read() (s string){
	outputs := [2]string{}
	outputs[0] = "exit"
	outputs[1] = "y"

	if (t.Iter < len(outputs)) {
		s = outputs[t.Iter]
		t.Iter++
	} else {
		s = "" // TODO make a proper ending character
	}

	return
}

func TestExitCommand(t *testing.T) {
	testMock := &TestInputOutput{0}
	game := game.GameManager{Out : testMock, In : testMock}
	game.Start()
	// Will hang if the game does not exit
}