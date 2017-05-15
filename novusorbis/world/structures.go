package world

type Structure interface {
	Mover
	Container
}

type BasicStructure struct {
	Mover
	Container
}

func NewStructure(name string, symbol string, position Position) (Structure) {
	return &BasicStructure{NewMover(name, symbol, position), NewContainer()}
}
