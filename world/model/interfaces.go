package model

type Mover interface {
	X() float64
	Y() float64
	Z() float64
	SetX(float64)
	SetY(float64)
	SetZ(float64)
}

type Namer interface {
	Name() string
	SetName(string)
}

type Describer interface {
	Description() string
	SetDescription(string)
}
