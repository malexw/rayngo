package rayngo

import (
	"image/color"
	"math"
	"github.com/malexw/vmath"
)

type Scene struct {
	Primitives []Primitive
	LightSrc Light
}

func NewScene() *Scene {
	s := Scene{
		Primitives: make([]Primitive, 0),
	}

	s.LightSrc = Light{vmath.Vec3{8.0, 7.0, 0.0}, ColorWhite(), 0.2, 0.7, 0.8}
	s.Primitives = append(s.Primitives, Primitive{Sphere{vmath.Vec3{1.0, 2.5, -13.0}, 2.5}, Material{Color{0, 0, 0.75, 1.0}, 50, false}})
	s.Primitives = append(s.Primitives, Primitive{Sphere{vmath.Vec3{3.0, 1.0, -7.0}, 1.0}, Material{Color{0, 0.5, 0, 1.0}, 20, false}})
	s.Primitives = append(s.Primitives, Primitive{Sphere{vmath.Vec3{-4.0, 2.0, -9.0}, 2.0}, Material{Color{0.75, 0, 0, 1.0}, 100, false}})
	s.Primitives = append(s.Primitives, Primitive{Plane{vmath.Vec3{0.0, 1.0, 0.0}, 0.0}, Material{Color{0.5, 0.5, 0.5, 1.0}, 10, true}})

	return &s
}

func (s *Scene) BackgroundColor(r Ray) color.RGBA {
	// TODO Correct for vectors that point straight up
	
	dir := r.Direction
	angle := (math.Atan(dir.Y/math.Abs(dir.Z))/math.Pi) + 0.5
	return color.RGBA{uint8(140 - (140 * angle)), uint8(208 - (208 * angle)), uint8(255 - (255 * angle)), 255}
}

func (s *Scene) IsShadowed(r Ray) bool {
	// For every object in the scene, check if the ray hits it.
	for _, prm := range s.Primitives {
		intersects, _, _, _ := prm.Geometry.RayCollision(r)
		if intersects {
			return true
		}
	}

	return false
}