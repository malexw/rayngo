package rayngo

import (
	"math"
)

type Material struct {
	Diffuse Color
	Specularity uint16
	Textured bool
}

func (m Material) Sample(u, v float64) Color {
	if m.Textured {
		if math.Mod((math.Ceil(u) + math.Ceil(v)), 2) == 0 {
			return Color{0.6, 0.6, 0.5, 1.0}
		}

		return Color{0.75, 0.75, 0.75, 1.0}
	}

	return m.Diffuse
}