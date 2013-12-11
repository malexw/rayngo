package rayngo

import (
	"github.com/malexw/vmath"
)

type Scene struct {
	Shapes []Shape
	LightSrc Light
}

func NewScene() *Scene {
	s := Scene{
		Shapes: make([]Shape, 1, 1),
	}

	s.LightSrc = Light{vmath.Vec3{8.0, 10.0, 5.0}}
	s.Shapes = append(s.Shapes, Shape{vmath.Vec3{3.0, 5.0, -10.0}, 1.0, vmath.Vec3{0, 128, 0}})
	s.Shapes = append(s.Shapes, Shape{vmath.Vec3{-7.0, 3.0, -15.0}, 2.0, vmath.Vec3{192, 192, 0}})

	return &s
}