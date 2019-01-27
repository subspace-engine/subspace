package model

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

type Passer interface {
	Passable() bool
	SetPassable(bool)
}

type Thing interface {
	Typer
	Namer
	Describer
	Passer
	Acter
	Shape
}
