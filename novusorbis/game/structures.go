package game

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