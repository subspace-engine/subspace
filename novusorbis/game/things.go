package game

type ThingStore struct {
	Things map[Position][]Thing
}

type Thing interface {
	Name()
	Symbol()
	Position()
}

type BasicThing struct {
	Name string
	Symbol string
	Position Position
}

func (store *ThingStore) Initialize() {
	store.Things = make(map[Position][]Thing)
}

func (store *ThingStore) AtPosition(p Position) (things []Thing, err error) {
	things, _ = store.Things[p]
	err = nil
	return
}

func (store *ThingStore) AddObjectAt(obj Thing, p Position) (err error) {
	const DEFAULT_STORE_SIZE = 3
	if store.Things[p] == nil {
		store.Things[p] = make([]Thing, 0, DEFAULT_STORE_SIZE)
	}
	store.Things[p] = append(store.Things[p], obj)
	err = nil
	return
}

func (thing *Thing) Name() (s string) {
	return thing.Name
}

func (thing *Thing) Symbol() (s string){
	return thing.Symbol
}
func (thing *Thing) Position() (p Position) {
	return thing.Position
}