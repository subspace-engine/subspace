package world

type NamedThing interface {
	Name() (s string)
}

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
	BasicNamedThing
	symbol string
	position Position
}

func NewNamedThing(Name string) (NamedThing){
	return &BasicNamedThing{name: Name}
}

func NewMover(Name string, Symbol string, Position Position) (Mover) {
	return &BasicMover{BasicNamedThing : NewNamedThing(Name), symbol: Symbol, position: Position}
}

type MoverStore interface {
	Initialize()
	AtPosition(p Position) (things []Mover, err error)
	AddObjectAt(obj Mover, p Position) (err error)
	MoveObjectTo(obj Mover, p Position) (err error)
	ShiftObjBy(obj Mover, p Position) (err error)
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

func (store *MapMoverStore) MoveObjectTo(obj Mover, p Position) (err error) {
	const DEFAULT_STORE_SIZE = 3
	if store.Movers[p] == nil {
		store.Movers[p] = make([]Mover, 0, DEFAULT_STORE_SIZE)
	}
	origPos := obj.Position()
	things, err := store.AtPosition(origPos)

	for index, element := range things {
		if(element == obj) {
			things = remove(things , index)
			store.Movers[origPos] = things
			break
		}
	}
	store.Movers[p] = append(store.Movers[p], obj)

	err = nil
	return
}

func (store *MapMoverStore) ShiftObjBy(obj Mover, p Position) (err error) {
	newPos, err := obj.Position().RelativePosition(p.X, p.Y, p.Z)
	store.MoveObjectTo(obj,newPos)
	return
}

func (thing *BasicMover) Name() (s string) {
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