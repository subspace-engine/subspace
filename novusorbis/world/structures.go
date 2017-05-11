package world

type Structure interface {
	Mover
}
type StructureStore interface {
	MoverStore
}

type BasicStructure struct {
	BasicMover
}

type MapStructureStore struct {
	MapMoverStore
}

func NewStructure() (Structure) {
	return &BasicStructure{}
}

func NewStructureStore() (StructureStore) {
	return &MapStructureStore{}
}