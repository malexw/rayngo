package rayngo

import (
	"github.com/malexw/vmath"
)

type Triangle struct {
	Vertices [3]vmath.Vec3
}

func (tr Triangle) RayCollision(r Ray) (bool, float64, vmath.Vec3, vmath.Vec3) {
	ab := tr.Vertices[1].Sub(tr.Vertices[0])
	ac := tr.Vertices[2].Sub(tr.Vertices[0])
	negDir := r.Direction.Scale(-1)

	normal := ab.Cross(ac)

	denom := negDir.Dot(normal)
	if denom <= 0.0 {
		return false, 0.0, vmath.Vec3{}, vmath.Vec3{}
	}

	ao := r.Origin.Sub(tr.Vertices[0])
	t := ao.Dot(normal)
	if t < 0.0 {
		return false, 0.0, vmath.Vec3{}, vmath.Vec3{}
	}

	e := negDir.Cross(ao)
	v := ac.Dot(e)
	if v < 0.0 || v > denom {
		return false, 0.0, vmath.Vec3{}, vmath.Vec3{}
	}
	w := -ab.Dot(e)
	if w < 0.0 || v + w > denom {
		return false, 0.0, vmath.Vec3{}, vmath.Vec3{}
	}

	denomInv := 1.0 / denom
	t = t * denomInv
	//v = v * denomInv
	//w = w * denomInv
	//u := 1.0 - v - w
	return true, t, r.Origin.Add(r.Direction.Scale(t)), normal.Normalize()
}

func (tr Triangle) Transform(m *vmath.Matrix4) Shape {
	return Triangle{[3]vmath.Vec3{m.Transform(tr.Vertices[0]),
								  m.Transform(tr.Vertices[1]),
								  m.Transform(tr.Vertices[2])}}
}