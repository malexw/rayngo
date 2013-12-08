package rayngo

import (
	// "math"
	"image/color"
	"github.com/malexw/vmath"
)

type Shape struct {
	Position vmath.Vec3
	Radius float64
	Color color.RGBA
}

func (s Shape) RayCollision(r Ray) bool {
	m := r.Origin.Sub(s.Position)
	b := m.Dot(r.Direction)
	c := m.Dot(m) - s.Radius * s.Radius
	if c > 0 && b > 0 {
		return false
	}
	disc := b*b - c
	if disc < 0 {
		return false
	}
	return true
	//t := -b - math.Sqrt(disc)
}