package rayngo

import (
	"github.com/malexw/vmath"
)

type Plane struct {
	Normal vmath.Vec3
	Distance float64
}

func (p Plane) RayCollision(r Ray) (bool, float64, vmath.Vec3, vmath.Vec3) {
	t := (p.Distance - p.Normal.Dot(r.Origin)) / p.Normal.Dot(r.Direction)

	if t >= 0.0 {
		intersection := r.Origin.Add(r.Direction.Scale(t))
		return true, t, intersection, p.Normal
	}

	return false, 0.0, vmath.Vec3{}, vmath.Vec3{}
}