package rayngo

import (
	"fmt"
	"github.com/malexw/vmath"
)

type Ray struct {
	Origin vmath.Vec3
	Direction vmath.Vec3
}

func (r Ray) String() string {
	return fmt.Sprintf("O: %s, D: %s", r.Origin.String(), r.Direction.String())
}