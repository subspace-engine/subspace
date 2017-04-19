package main

import (
	"github.com/subspace-engine/subspace/novaterram/ui"
	"github.com/subspace-engine/subspace/novaterram/game"
)

func main() {
	inOut := ui.NewInputOutput()
	game := game.GameManager{Out : inOut, In : inOut}
	game.Start()
}
