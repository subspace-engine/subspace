package main

import (
	"github.com/subspace-engine/subspace/simtext/ui"
	"github.com/subspace-engine/subspace/simtext/game"
)

func main() {
	inOut := ui.NewInputOutput()
	game := game.LoopHandler{Out : inOut, In : inOut}
	game.Start()
}
