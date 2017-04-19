package main

import (
	"github.com/subspace-engine/subspace/xenoterra/ui"
	"github.com/subspace-engine/subspace/xenoterra/game"
)

func main() {
	inOut := ui.NewInputOutput()
	game := game.GameManager{Out : inOut, In : inOut}
	game.Start()
}
