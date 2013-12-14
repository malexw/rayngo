package rayngo

import (
	"github.com/malexw/vmath"
)

type Light struct {
	Position vmath.Vec3
	Diffuse Color
	AmbientCoeff, DiffuseCoeff, SpecularCoeff float64
}