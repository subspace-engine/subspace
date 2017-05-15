package world

type Structure interface {
	Mover
	Container
}

type BasicStructure struct {
	Mover
	Container
}

type Container interface {
	AddObject(thing NamedThing)
	GetContents() (things []NamedThing)
}

type BasicContainer struct {
	things []NamedThing
}

func NewStructure(name string, symbol string, position Position) (Structure) {
	return &BasicStructure{NewMover(name, symbol, position), NewContainer()}
}

func NewContainer() (Container) {
	return &BasicContainer{things : make([]NamedThing, 0, 4)}
}

func (c *BasicContainer) AddObject(thing NamedThing) {
	c.things = append(c.things, thing)
}

func (c *BasicContainer) GetContents() ([]NamedThing){
	return c.things
}
