package main

import (
	"github.com/subspace-engine/subspace/novusorbis/game"
	"github.com/subspace-engine/subspace/novusorbis/ui"
)

func main() {
	inOut := ui.NewInputOutput()
	game := game.GameManager{Out: inOut, In: inOut}
	game.Start()
}
