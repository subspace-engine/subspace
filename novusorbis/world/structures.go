package world

type Structure interface {
	Thing
}
type StructureStore interface {
	ThingStore
}

type BasicStructure struct {
	BasicThing
}

type MapStructureStore struct {
	MapThingStore
}

func NewStructure() (Structure) {
	return &BasicStructure{}
}

func NewStructureStore() (StructureStore) {
	return &MapStructureStore{}
}