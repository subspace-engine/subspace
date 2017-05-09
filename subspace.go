package main

import (
	"github.com/subspace-engine/subspace/novusorbis/ui"
	"github.com/subspace-engine/subspace/novusorbis/game"
)

func main() {
	inOut := ui.NewInputOutput()
	questionAsker := game.QuestionAsker{Out : inOut, In : inOut}
	baseFactory := game.BaseFactory{questionAsker}
	game := game.GameManager{Out : inOut, In : inOut, BaseFactory: baseFactory}
	game.Start()
}
