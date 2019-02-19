package model

import "github.com/subspace-engine/subspace/util"

type Identity interface {
	ID() string
	SetID(string)
}

type Namer interface {
	Name() string
	SetName(string)
}

type Describer interface {
	Description() string
	SetDescription(string)
}

type Typer interface {
	Type() string
	SetType(string)
}

type Mover interface {
	Move(Thing, util.Vec3) bool
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

type Sayer interface { // for message handling and reporting of events
	Say(text string)
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
	Sayer
	Actions() Actor
}

type MobileThing interface {
	Thing
	StepSize() float64
	SetStepSize(float64)
	Direction() float64
	SetDirection(float64)
	Identity
}
