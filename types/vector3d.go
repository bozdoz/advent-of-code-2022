package types

import (
	"fmt"
	"math"
)

type Vector3d struct {
	X, Y, Z int
}

func NewVector3d(x, y, z int) Vector3d {
	return Vector3d{x, y, z}
}

func (v *Vector3d) LengthSquared() int {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3d) Length() float64 {
	return math.Sqrt(float64(v.LengthSquared()))
}

func (v Vector3d) Add(b Vector3d) Vector3d {
	return Vector3d{
		v.X + b.X,
		v.Y + b.Y,
		v.Z + b.Z,
	}
}

func (v *Vector3d) Subtract(b Vector3d) Vector3d {
	return Vector3d{
		v.X - b.X,
		v.Y - b.Y,
		v.Z - b.Z,
	}
}

func (v *Vector3d) Divide(b Vector3d) Vector3d {
	var x, y, z int
	if b.X == 0 || v.X == 0 {
		x = 0
	} else {
		x = v.X / b.X
	}
	if b.Y == 0 || v.Y == 0 {
		y = 0
	} else {
		y = v.Y / b.Y
	}
	if b.Z == 0 || v.Z == 0 {
		z = 0
	} else {
		z = v.Z / b.Z
	}

	return Vector3d{x, y, z}
}

func (v *Vector3d) Multiply(b Vector3d) Vector3d {
	return Vector3d{
		v.X * b.X,
		v.Y * b.Y,
		v.Z * b.Z,
	}
}

func (v *Vector3d) DistanceTo(b Vector3d) float64 {
	return v.Subtract(b).Length()
}

func (v *Vector3d) String() string {
	return fmt.Sprint(*v)
}
