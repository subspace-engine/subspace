package world

type Thing interface {
	Name() (s string)
	Symbol() (s string)
	Position() (p Position)
	SetPosition(p Position)
}

type BasicThing struct {
	name string
	symbol string
	position Position
}

func NewThing(Name string, Symbol string, Position Position) (Thing) {
	return &BasicThing{name: Name, symbol: Symbol, position: Position}
}


type ThingStore interface {
	Initialize()
	AtPosition(p Position) (things []Thing, err error)
	AddObjectAt(obj Thing, p Position) (err error)
	MoveObjectTo(obj Thing, p Position) (err error)
	ShiftObjBy(obj Thing, p Position) (err error)
}

type MapThingStore struct {
	Things map[Position][]Thing
}

func (store *MapThingStore) Initialize() {
	store.Things = make(map[Position][]Thing)
}

func (store *MapThingStore) AtPosition(p Position) (things []Thing, err error) {
	things, _ = store.Things[p]
	err = nil
	return
}

func (store *MapThingStore) AddObjectAt(obj Thing, p Position) (err error) {
	const DEFAULT_STORE_SIZE = 3
	if store.Things[p] == nil {
		store.Things[p] = make([]Thing, 0, DEFAULT_STORE_SIZE)
	}
	store.Things[p] = append(store.Things[p], obj)
	err = nil
	return
}

func remove(s []Thing, i int) []Thing {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (store *MapThingStore) MoveObjectTo(obj Thing, p Position) (err error) {
	const DEFAULT_STORE_SIZE = 3
	if store.Things[p] == nil {
		store.Things[p] = make([]Thing, 0, DEFAULT_STORE_SIZE)
	}
	origPos := obj.Position()
	things, err := store.AtPosition(origPos)

	for index, element := range things {
		if(element == obj) {
			things = remove(things , index)
			store.Things[origPos] = things
			break
		}
	}
	store.Things[p] = append(store.Things[p], obj)

	err = nil
	return
}

func (store *MapThingStore) ShiftObjBy(obj Thing, p Position) (err error) {
	newPos, err := obj.Position().RelativePosition(p.X, p.Y, p.Z)
	store.MoveObjectTo(obj,newPos)
	return
}

func (thing *BasicThing) Name() (s string) {
	return thing.name
}

func (thing *BasicThing) Symbol() (s string){
	return thing.symbol
}

func (thing *BasicThing) Position() (p Position) {
	return thing.position
}

func (thing *BasicThing) SetPosition(p Position) {
	thing.position = p
}