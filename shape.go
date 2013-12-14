package rayngo

import (
	"github.com/malexw/vmath"
)

type Shape interface {
	RayCollision(Ray) (bool, float64, vmath.Vec3, vmath.Vec3)
}
