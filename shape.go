package rayngo

import (
	"math"
	//"image/color"
	"github.com/malexw/vmath"
)

type Shape struct {
	Position vmath.Vec3
	Radius float64
}

func (s Shape) RayCollision(r Ray) (bool, vmath.Vec3) {
	m := r.Origin.Sub(s.Position)
	b := m.Dot(r.Direction)
	c := m.Dot(m) - s.Radius * s.Radius
	if c > 0 && b > 0 {
		return false, vmath.Vec3{}
	}
	disc := b*b - c
	if disc < 0 {
		return false, vmath.Vec3{}
	}
	t := -b - math.Sqrt(disc)
	intersection := r.Origin.Add(r.Direction.Scale(t))
	return true, intersection
	//t := -b - math.Sqrt(disc)
}