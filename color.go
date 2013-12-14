package rayngo

import (
	"image/color"
	"math"
)

type Color struct {
	R, G, B, A float64
}

func ColorBlack() Color {
	return Color{0.0, 0.0, 0.0, 1.0}
}

func ColorWhite() Color {
	return Color{1.0, 1.0, 1.0, 1.0}
}

func (c Color) Add(rhs Color) Color {
	return Color{c.R+rhs.R, c.G+rhs.G, c.B+rhs.B, c.A}
}

func (c Color) Attenuate(n float64) Color {
	return Color{c.R*n, c.G*n, c.B*n, c.A}
}

func (c Color) ToImageColor() color.RGBA {
	red := uint8(clamp(0, c.R, 1.0)*255)
	green := uint8(clamp(0, c.G, 1.0)*255)
	blue := uint8(clamp(0, c.B, 1.0)*255)
	alpha := uint8(clamp(0, c.A, 1.0)*255)
	return color.RGBA{red, green, blue, alpha}
}

func clamp(min, val, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}