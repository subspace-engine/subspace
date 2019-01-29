package model

import (
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

type BasicThing struct {
	objType int
	Shape
	name        string
	description string
	passable    bool
	Actor
	Locater
}

func MakeBasicThing(name string, description string) *BasicThing {
	return &BasicThing{0, &Point{util.Vec3{0, 0, 0}}, name, description, false, MakeActionManager(), makeLocater()}
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

func (self BasicThing) Move(obj Thing, pos util.Vec3) int {
	return 1 // not implemented for basic thing
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

type MobileThing struct {
	*BasicThing
	stepSize  float64
	direction float64
}

func MakeMobileThing(name string, description string) *MobileThing {
	t := &MobileThing{MakeBasicThing(name, description), 1, 0}
	t.RegisterAction("move", func(action Action) int {
		pos := util.Vec3{math.Sin(t.direction) * t.stepSize,
			0,
			-math.Cos(t.direction) * t.stepSize}
		if t.Location() != nil {
			t.Location().Move(t, pos)
			return 1
		}
		return 0
	})
	return t
}

func (self MobileThing) StepSize() float64 {
	return self.stepSize
}

func (self *MobileThing) SetStepSize(size float64) {
	self.stepSize = size
}

func (self MobileThing) Direction() float64 {
	return self.direction
}

func (self *MobileThing) SetDirection(direction float64) {
	self.direction = direction
}
