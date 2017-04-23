package world

import "io"

type Saver interface {
	save(io.Writer) bool
}

type Loader interface {
	load(io.Reader) bool
}

type Object map[string](interface{})

func NewObject() Object {
	obj := make(map[string](interface{}))
	return obj
}

type Shape interface {
	MaxX() float64
	MaxY() float64
	MaxZ() float64
	MinX() float64
	MinY() float64
	MinZ() float64
	IsTouching(othr Shape) bool
}

type Box struct {
	x      float64
	y      float64
	z      float64
	width  float64
	height float64
	depth  float64
}

func NewBox(x float64, y float64, z float64, width float64, height float64, depth float64) Box {
	return Box{x, y, z, width, height, depth}
}

func (b Box) MinX() float64 {
	return b.x
}

func (b Box) MinY() float64 {
	return b.y
}

func (b Box) MinZ() float64 {
	return b.z
}

func (b Box) MaxX() float64 {
	return b.x + b.width
}

func (b Box) MaxY() float64 {
	return b.y + b.height
}

func (b Box) MaxZ() float64 {
	return b.z + b.depth
}

func (b Box) IsTouching(othr Shape) bool {
	return !((b.MinX() > othr.MaxX() || b.MaxX() < othr.MinX()) && (b.MinY() > othr.MaxY() || b.MaxY() < othr.MinY()) && (b.MinZ() > othr.MaxZ() || b.MaxZ() < othr.MinZ()))
}

type WObject struct {
	Shape
	Object
	Relations map[string](*WObject)
}

func NewWObject(name string, shape Shape) WObject {
	obj := WObject{shape, NewObject(), make(map[string](*WObject))}
	obj.SetName(name)
	return obj
}

func (obj WObject) RelatedTo(child *WObject) bool {
	for _, val := range obj.Relations {
		if val == child {
			return true
		}
	}
	return false
}

func (obj Object) Name() string {
	return obj["name"].(string)
}

func (obj Object) SetName(name string) {
	obj["name"] = name
}

func (obj WObject) PutIn(container WObject) {
	obj.Relations["in"] = &container
	container.Relations["contains"] = &obj
}
