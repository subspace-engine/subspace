package model

import (
	"github.com/subspace-engine/subspace/util"
)

type Shape interface {
	Position() util.Vec3
	SetPosition(util.Vec3)
	Box() util.Vec3
	IsTouching(Shape) bool
}

type Point struct {
	pos util.Vec3
}

func MakePoint(v util.Vec3) *Point {
	return &Point{v}
}

func (p Point) Position() util.Vec3 {
	return p.pos
}

func (p *Point) SetPosition(v util.Vec3) {
	p.pos = v
}

func (p Point) Box() util.Vec3 {
	return util.Vec3{0, 0, 0}
}

func (p *Point) IsTouching(othr Shape) bool {
	switch t := othr.(type) {
	case *Point:
		return p.pos.Equals(t.Position())
	default:
		return othr.IsTouching(p)
	}
}
