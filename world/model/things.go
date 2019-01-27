package model

import "github.com/subspace-engine/subspace/util"

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
	actions     *ActionManager
	Locater
}

func MakeBasicThing(name string, description string) Thing {
	return &BasicThing{0, &Point{util.Vec3{0, 0, 0}}, name, description, false, MakeActionManager(), makeLocater()}
}

func MakePassableThing(name string, description string, passable bool) Thing {
	t := MakeBasicThing(name, description)
	t.SetPassable(passable)
	return t
}

func MakeTypedThing(objType int, name string, description string, passable bool) Thing {
	t := MakePassableThing(name, description, passable)
	t.SetType(objType)
	return t
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

func (self *BasicThing) Act(action Action) {
	self.actions.Act(action)
}

func (self *BasicThing) RegisterAction(tag string, response func(Action) int) {
	self.actions.RegisterAction(tag, response)
}

func (thing BasicThing) IsRoot() bool {
	return false
}
