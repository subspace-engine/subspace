package main

import (
	"github.com/subspace-engine/subspace/simtext"
	"github.com/subspace-engine/subspace/ui"
)

func main() {
	inOut := ui.NewInputOutput()
	game := simtext.Game{Out : inOut, In : inOut}
	game.Start()
}
