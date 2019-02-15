package model

import (
	"github.com/google/uuid"
	"github.com/subspace-engine/subspace/util"
	"math"
)

type basicLocater struct {
	location Thing
	children []Thing
}

func makeLocater() Locater {
	return &basicLocater{nil, make([]Thing, 0, 0)}
}

func (loc basicLocater) Location() Thing {
	return loc.location
}

func (loc *basicLocater) SetLocation(thing Thing) {
	loc.location = thing
}

func (loc basicLocater) Children() []Thing {
	return loc.children
}

func (loc *basicLocater) AddChild(child Thing) {
	loc.children = append(loc.children, child)
}

func (loc *basicLocater) RemoveChild(child Thing) {
	if len(loc.children) == 0 {
		return
	}
	indx := 0
	for i, val := range loc.children {
		if val == child {
			indx = i
		}
	}
	loc.children[indx] = loc.children[len(loc.children)-1]
	loc.children = loc.children[:len(loc.children)-1]
}

type UuidIdentity struct {
	id string
}

func (self *UuidIdentity) SetID(id string) {
	self.id = id
}

func (self UuidIdentity) ID() string {
	return self.id
}

func MakeIdentity() Identity {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic("Unable to make uuid.")
	}
	return &UuidIdentity{uuid.String()}
}

type BasicThing struct {
	objType int
	Shape
	name        string
	description string
	passable    bool
	Actor
	Locater
	Identity
}

func MakeBasicThing(name string, description string) *BasicThing {
	return &BasicThing{0, &Point{util.Vec3{0, 0, 0}}, name, description, false, MakeActionManager(), makeLocater(), MakeIdentity()}
}

func MakePassableThing(name string, description string, passable bool) *BasicThing {
	t := MakeBasicThing(name, description)
	t.SetPassable(passable)
	return t
}

func MakeTypedThing(objType int, name string, description string, passable bool) *BasicThing {
	t := MakePassableThing(name, description, passable)
	t.SetType(objType)
	return t
}

func (self BasicThing) Move(obj Thing, pos util.Vec3) bool {
	return false // not implemented for basic thing
}

func (self *BasicThing) Name() string {
	return self.name
}

func (self *BasicThing) SetName(name string) {
	self.name = name
}

func (self *BasicThing) Description() string {
	return self.description
}

func (self *BasicThing) SetDescription(description string) {
	self.description = description
}

func (self *BasicThing) Passable() bool {
	return self.passable
}

func (self *BasicThing) SetPassable(passable bool) {
	self.passable = passable
}

func (self *BasicThing) Type() int {
	return self.objType
}

func (self *BasicThing) SetType(objType int) {
	self.objType = objType
}

func (thing BasicThing) IsRoot() bool {
	return false
}

func (self *BasicThing) Say(text string) {
	// Send it up to the world to handle
	if self.Location() != nil {
		self.Location().Say(text)
	}
}

func (self *BasicThing) Actions() Actor {
	return self.Actor
}

type BasicMobileThing struct {
	*BasicThing
	stepSize  float64
	direction float64
}

func MakeMobileThing(name string, description string) MobileThing {
	t := &BasicMobileThing{MakeBasicThing(name, description), 1, 0}
	t.RegisterAction("forward", func(action Action) bool {
		if action.Source == nil {
			return false
		}
		thing, ok := action.Source.(MobileThing)
		if !ok {
			return false
		}

		pos := util.Vec3{math.Sin(t.direction)*t.stepSize + 0.000001,
			0,
			-math.Cos(t.direction)*t.stepSize + 0.000001}
		if thing.Location() != nil {
			return thing.Location().Move(thing, pos)
		}
		return false
	})
	t.RegisterAction("sidestep left", func(action Action) bool {
		if action.Source == nil {
			return false
		}
		thing, ok := action.Source.(MobileThing)
		if !ok {
			return false
		}

		pos := util.Vec3{math.Sin(t.direction-(math.Pi/2))*t.stepSize/2.0 + 0.000001,
			0,
			-math.Cos(t.direction-(math.Pi/2))*t.stepSize/2.0 + 0.000001}
		if thing.Location() != nil {
			return thing.Location().Move(thing, pos)
		}
		return false
	})
	t.RegisterAction("sidestep right", func(action Action) bool {
		if action.Source == nil {
			return false
		}
		thing, ok := action.Source.(MobileThing)
		if !ok {
			return false
		}

		pos := util.Vec3{math.Sin(t.direction+(math.Pi/2))*t.stepSize/2.0 + 0.000001,
			0,
			-math.Cos(t.direction+(math.Pi/2))*t.stepSize/2.0 + 0.000001}
		if thing.Location() != nil {
			return thing.Location().Move(thing, pos)
		}
		return false
	})
	t.RegisterAction("reverse", func(action Action) bool {
		if action.Source == nil {
			return false
		}
		thing, ok := action.Source.(MobileThing)
		if !ok {
			return false
		}

		pos := util.Vec3{math.Sin(t.direction+math.Pi)*t.stepSize/2.0 + 0.000001,
			0,
			-math.Cos(t.direction+math.Pi)*t.stepSize/2.0 + 0.000001}
		if thing.Location() != nil {
			return thing.Location().Move(thing, pos)
		}
		return false
	})
	t.RegisterAction("turn left", func(action Action) bool {
		source, ok := action.Source.(MobileThing)
		if !ok {
			return false
		}
		source.SetDirection(math.Mod(source.Direction()-(math.Pi/2.0)+(2*math.Pi), 2*math.Pi))
		return true
	})
	t.RegisterAction("turn right", func(action Action) bool {
		source, ok := action.Source.(MobileThing)
		if !ok {
			return false
		}
		source.SetDirection(math.Mod(source.Direction()+(math.Pi/2.0)+(2*math.Pi), 2*math.Pi))
		return true
	})

	return t
}

func (self BasicMobileThing) StepSize() float64 {
	return self.stepSize
}

func (self *BasicMobileThing) SetStepSize(size float64) {
	self.stepSize = size
}

func (self BasicMobileThing) Direction() float64 {
	return self.direction
}

func (self *BasicMobileThing) SetDirection(direction float64) {
	self.direction = direction
}
