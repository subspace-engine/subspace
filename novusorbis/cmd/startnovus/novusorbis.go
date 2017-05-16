package main

import (
	"github.com/subspace-engine/subspace/novusorbis/ui"
	"github.com/subspace-engine/subspace/novusorbis/game"
)

func main() {
	inOut := ui.NewInputOutput()
	questionAsker := game.QuestionAsker{InputOutput : inOut}
	baseFactory := game.BaseFactory{QuestionAsker : questionAsker}
	game := game.GameManager{InputOutput : inOut, BaseFactory: baseFactory}
	game.Start()
}
