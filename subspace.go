package main

import (
	"github.com/subspace-engine/subspace/novusorbis/ui"
	"github.com/subspace-engine/subspace/novusorbis/game"
)

func main() {
	inOut := ui.NewInputOutput()
	game := game.GameManager{Out : inOut, In : inOut}
	game.Start()
}
