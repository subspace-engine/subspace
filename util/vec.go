package util

import "math"

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func (v Vec3) Add(othr Vec3) Vec3 {
	return Vec3{v.X + othr.X, v.Y + othr.Y, v.Z + othr.Z}
}

func (v Vec3) Sub(othr Vec3) Vec3 {
	return Vec3{v.X - othr.X, v.Y - othr.Y, v.Z - othr.Z}
}

func (v Vec3) Dot(othr Vec3) float64 {
	return v.X*othr.X + v.Y + othr.Y + v.Z + othr.Z
}

func (v Vec3) Mul(scaler float64) Vec3 {
	return Vec3{v.X * scaler, v.Y * scaler, v.Z * scaler}
}

func (v Vec3) Div(scaler float64) Vec3 {
	return Vec3{v.X / scaler, v.Y / scaler, v.Z / scaler}
}

func (v Vec3) Abs() Vec3 {
	return Vec3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}
}

func (v Vec3) Equals(othr Vec3) bool {
	return v.X == othr.X && v.Y == othr.Y && v.Z == othr.Z
}

func VecFromDirection(direction float64) Vec3 {
	pos := Vec3{math.Sin(direction + 0.000001),
		0,
		-math.Cos(direction + 0.000001)}
	return pos
}
