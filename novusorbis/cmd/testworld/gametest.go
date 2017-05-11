package main

/**
 * This file sets up a small test world, skipping any intro etc.
 */

import (
	"github.com/subspace-engine/subspace/novusorbis/ui"
	"github.com/subspace-engine/subspace/novusorbis/game"
	"github.com/subspace-engine/subspace/novusorbis/world"
)

type testGameManager struct {
	*game.GameManager
}

func main() {
	inOut := ui.NewInputOutput()
	questionAsker := game.QuestionAsker{InputOutput : inOut}
	baseFactory := game.BaseFactory{QuestionAsker : questionAsker}
	g := game.GameManager{InputOutput : inOut, BaseFactory: baseFactory}
	t := testGameManager{&g}
	t.initializeTestWorld()
	t.InitializeCommandsMap()
	t.SetUpDirectionMaps()
	t.MainLoop()
}

func (g *testGameManager) initializeTestWorld() {
	w := &world.World{}
	w.Size = 5
	g.World = w
	mid := w.Size/2

	c := world.Position{mid,mid,mid}
	w.Cursor = c

	w.GenerateTerrain()

	thingStore := &world.MapMoverStore{}
	thingStore.Initialize()
	w.Things = thingStore

	w.Structures = world.NewStructureStore()

	w.MainColonist = g.CreateDefaultColonist() // TODO
	w.MainBase = g.CreateDefaultBase() // TODO

	pos := world.Position{mid,mid,mid}
	mainAvatar := g.World.MainColonist.Avatar
	mainAvatar.SetPosition(pos)
	thingStore.AddObjectAt(mainAvatar, pos)

	baseAvatar := g.World.MainBase.Avatar
	baseAvatar.SetPosition(pos)
	thingStore.AddObject(baseAvatar)

	icePos, _ := world.NewPosition(mid).RelativePosition(-1,-1,0)
	ice := world.NewMover("ice", "i", icePos)
	thingStore.AddObject(ice)

	return
}