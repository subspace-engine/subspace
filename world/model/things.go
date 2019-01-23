package model

type MoveObj struct {
	x float64
	y float64
	z float64
}

func (self *MoveObj) X() float64 {
	return self.x
}

func (self *MoveObj) Y() float64 {
	return self.y
}

func (self *MoveObj) Z() float64 {
	return self.z
}

func (self *MoveObj) SetX(x float64) {
	self.x = x
}

func (self *MoveObj) SetY(y float64) {
	self.y = y
}

func (self *MoveObj) SetZ(z float64) {
	self.z = z
}

type BasicThing struct {
	objType int
	Mover
	name        string
	description string
	passable    bool
	actions     *ActionManager
}

func MakeBasicThing(name string, description string) Thing {
	return &BasicThing{0, &MoveObj{0, 0, 0}, name, description, false, MakeActionManager()}
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
