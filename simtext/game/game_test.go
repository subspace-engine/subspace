package game_test

import (
	"testing"
	"github.com/subspace-engine/subspace/simtext/game"
)

type InputOutput interface {
	Print(s string)
	Println(s string)
	Read() (s string)
}


type TestInputOutput struct{}

func (t *TestInputOutput) Print(s string) {
	// Do nothing
}

func (t *TestInputOutput) Println(s string) {
	// Do nothing
}

func (t *TestInputOutput) Read() (s string){
	s = "exit"
	return
}

func TestMainloop(t *testing.T) {
	testMock := &TestInputOutput{}
	game := game.LoopHandler{Out : testMock, In : testMock}
	game.Start()
	if(!true) {
		t.Error("Not true")
	}
}