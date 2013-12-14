package rayngo

import (
	"math"
	"github.com/malexw/vmath"
)

type Sphere struct {
	Position vmath.Vec3
	Radius float64
}

func (s Sphere) RayCollision(r Ray) (bool, float64, vmath.Vec3, vmath.Vec3) {
	m := r.Origin.Sub(s.Position)
	b := m.Dot(r.Direction)
	c := m.Dot(m) - s.Radius * s.Radius
	if c > 0 && b > 0 {
		return false, 0.0, vmath.Vec3{}, vmath.Vec3{}
	}
	disc := b*b - c
	if disc < 0 {
		return false, 0.0, vmath.Vec3{}, vmath.Vec3{}
	}
	t := -b - math.Sqrt(disc)
	intersection := r.Origin.Add(r.Direction.Scale(t))
	normal := (intersection.Sub(s.Position)).Normalize()
	return true, t, intersection, normal
}