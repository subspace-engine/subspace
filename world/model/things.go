package model


type MoveObj struct {
	x float64
	y float64
	z float64
}

func    (self *MoveObj)X() float64 {
	return self.x
}

func (self *MoveObj)Y() float64 {
	return self.y
}

func (self *MoveObj)Z() float64 {
	return self.z
}

func (self*MoveObj)SetX(x float64) {
	self.x=x
}

func (self*MoveObj)SetY(y float64) {
	self.y=y;
}

func (self*MoveObj)SetZ(z float64) {
	self.z=z
}


type Thing struct {
	objType int
	MoveObj
	name string
	description string
	passable bool
	actions *ActionManager
}

func MakeBasicThing(name string, description string) *Thing {
	return &Thing {0, MoveObj{0,0,0}, name, description, false, MakeActionManager()}
}

func MakePassableThing(name string, description string, passable bool) *Thing {
	t:=MakeBasicThing(name, description)
	t.SetPassable(passable)
	return t
}

func MakeTypedThing(objType int, name string, description string, passable bool) *Thing {
	t:=MakePassableThing(name,description,passable)
	t.SetType(objType)
	return t
}

func (self*Thing)Name() string {
	return self.name
}

func (self*Thing)SetName(name string) {
	self.name=name
}

func (self*Thing)Description() string {
	return self.description
}

func (self*Thing)SetDescription(description string) {
	self.description=description
	}


func (self*Thing)Passable() bool {
	return self.passable
}

func (self*Thing)SetPassable(passable bool) {
	self.passable=passable
}

func (self*Thing)Type() int {
	return self.objType
}

func (self*Thing)SetType(objType int) {
	self.objType=objType
}


func (self*Thing)Act(action Action) {
	self.actions.Act(action)
}

func (self *Thing)RegisterAction(tag string, response func(Action)int) {
	self.actions.RegisterAction(tag, response)
}
