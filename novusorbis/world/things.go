package world

type NamedThing interface {
	Name() (s string)
}

// For now, a mover is anything with a position, even if it doesn't actually move
type Mover interface {
	NamedThing
	Symbol() (s string)
	Position() (p Position)
	SetPosition(p Position)
}

type BasicNamedThing struct {
	name string
}

type BasicMover struct {
	NamedThing
	symbol string
	position Position
}

func NewNamedThing(Name string) (NamedThing){
	return &BasicNamedThing{name: Name}
}

func NewMover(Name string, Symbol string, Position Position) (Mover) {
	return &BasicMover{NamedThing : NewNamedThing(Name), symbol: Symbol, position: Position}
}

type MoverStore interface {
	Initialize()
	AtPosition(p Position) (things []Mover, err error)
	AddObjectAt(obj Mover, p Position) (err error)
	MoveObjectTo(obj Mover, p Position) (err error)
	ShiftObjBy(obj Mover, p Position) (err error)
	Remove(obj Mover, p Position) (err error)
}

type MapMoverStore struct {
	Movers map[Position][]Mover
}

func (store *MapMoverStore) Initialize() {
	store.Movers = make(map[Position][]Mover)
}

func (store *MapMoverStore) AtPosition(p Position) (movers []Mover, err error) {
	movers, _ = store.Movers[p]
	err = nil
	return
}

func (store *MapMoverStore) AddObjectAt(obj Mover, p Position) (err error) {
	const DEFAULT_STORE_SIZE = 3
	if store.Movers[p] == nil {
		store.Movers[p] = make([]Mover, 0, DEFAULT_STORE_SIZE)
	}
	store.Movers[p] = append(store.Movers[p], obj)
	err = nil
	return
}

func (store *MapMoverStore) AddObject(obj Mover) (err error) {
	const DEFAULT_STORE_SIZE = 3
	p := obj.Position()
	if store.Movers[p] == nil {
		store.Movers[p] = make([]Mover, 0, DEFAULT_STORE_SIZE)
	}
	store.Movers[p] = append(store.Movers[p], obj)
	err = nil
	return
}

func remove(s []Mover, i int) []Mover {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (store *MapMoverStore) Remove(obj Mover, p Position) (err error) {
	things, err := store.AtPosition(p)
	for index, element := range things {
		if(element == obj) {
			things = remove(things , index)
			store.Movers[p] = things
			return
		}
	}
	return
}

func (store *MapMoverStore) MoveObjectTo(obj Mover, p Position) (err error) {
	const DEFAULT_STORE_SIZE = 3
	if store.Movers[p] == nil {
		store.Movers[p] = make([]Mover, 0, DEFAULT_STORE_SIZE)
	}
	origPos := obj.Position()
	store.Remove(obj, origPos)
	store.Movers[p] = append(store.Movers[p], obj)

	err = nil
	return
}

func (store *MapMoverStore) ShiftObjBy(obj Mover, p Position) (err error) {
	newPos, err := obj.Position().RelativePosition(p.X, p.Y, p.Z)
	store.MoveObjectTo(obj,newPos)
	return
}

func (thing *BasicNamedThing) Name() (s string) {
	return thing.name
}

func (thing *BasicMover) Symbol() (s string){
	return thing.symbol
}

func (thing *BasicMover) Position() (p Position) {
	return thing.position
}

func (thing *BasicMover) SetPosition(p Position) {
	thing.position = p
}


type Container interface {
	AddObject(thing NamedThing)
	GetContents() (things []NamedThing)
}

type BasicContainer struct {
	things []NamedThing
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
