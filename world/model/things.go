package model

type MoveObj struct {
	x float64
	y float64
	z float64
}

func (self *MoveObj)X() float64 {
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
	MoveObj
	name string
	description string
}

func MakeThing(name string, description string) Thing {
	return Thing {MoveObj{0,0,0}, name, description}
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
