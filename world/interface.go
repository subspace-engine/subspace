package world

import (
	engine "github.com/subspace-engine/subspace/engine"
)

type PlayerWorld interface {
	Actor() chan<- engine.PlayerAction
	Updater(id int) <-chan *engine.ServerUpdate
	Run(chan struct{})
	RequestPlayer() Player //will change
}

type Player interface {
	Id() int
}
