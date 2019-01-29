package model

import "github.com/subspace-engine/subspace/util"

type Namer interface {
	Name() string
	SetName(string)
}

type Describer interface {
	Description() string
	SetDescription(string)
}

type Typer interface {
	Type() int
	SetType(int)
}

type Mover interface {
	Move(Thing, util.Vec3) int
}

type Passer interface {
	Passable() bool
	SetPassable(bool)
}

type Locater interface {
	Location() Thing
	SetLocation(Thing)
	Children() []Thing
	AddChild(Thing)
	RemoveChild(Thing)
}

type Thing interface {
	Typer
	Namer
	Describer
	Passer
	Actor
	Locater
	Mover
	Shape
	IsRoot() bool
}
