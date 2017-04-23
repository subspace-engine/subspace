package world

type Point2d struct {
	x float64
	y float64
}

func (point Point2d) X() float64 {
	return point.x
}

func (point *Point2d) SetX(x float64) {
	point.x = x
}

func (point Point2d) Y() float64 {
	return point.y
}

func (point *Point2d) SetY(y float64) {
	point.y = y
}

func (point Point2d) Z() {
	return 0
}

func (point *Point3d) SetZ() {
}

type Point3d struct {
	Point2d
	z float64
}

func (point Point3d) Z() float64 {
	return point.z
}

func (point *Point3d) SetZ(z float64) {
	point.z = z
}
