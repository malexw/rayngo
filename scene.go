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
	s.Primitives = append(s.Primitives, Primitive{Sphere{vmath.Vec3{-4.0, 1.5, -9.0}, 1.5}, Material{Color{0.75, 0, 0, 1.0}, 100, false}})
	s.Primitives = append(s.Primitives, Primitive{Plane{vmath.Vec3{0.0, 1.0, 0.0}, 0.0}, Material{Color{0.5, 0.5, 0.5, 1.0}, 10, true}})
	
	pyrMat := Material{Color{1.0, 0.6, 0.2, 1.0}, 200, false}
	tri1 := Triangle{[3]vmath.Vec3{vmath.Vec3{-8.0, 0.0, -18.0}, vmath.Vec3{-8.0, 8.0, -21.0}, vmath.Vec3{-12.0, 0.0, -24.0}}}
	s.Primitives = append(s.Primitives, Primitive{tri1, pyrMat})

	tri2 := Triangle{[3]vmath.Vec3{vmath.Vec3{-4.0, 0.0, -24.0}, vmath.Vec3{-8.0, 8.0, -21.0}, vmath.Vec3{-8.0, 0.0, -18.0}}}
	s.Primitives = append(s.Primitives, Primitive{tri2, pyrMat})

	tri3 := Triangle{[3]vmath.Vec3{vmath.Vec3{-12.0, 0.0, -24.0}, vmath.Vec3{-8.0, 8.0, -21.0}, vmath.Vec3{-4.0, 0.0, -24.0}}}
	s.Primitives = append(s.Primitives, Primitive{tri3, pyrMat})

	return &s
}

func (s *Scene) BackgroundColor(r Ray) color.RGBA {
	// TODO Correct for vectors that point straight up
	
	dir := r.Direction
	angle := (math.Atan(dir.Y/math.Abs(dir.Z))/math.Pi) + 0.5
	return color.RGBA{uint8(140 - (140 * angle)), uint8(208 - (208 * angle)), uint8(255 - (255 * angle)), 255}
}

func (s *Scene) IsShadowed(r Ray) bool {
	dtl := r.Origin.Sub(s.LightSrc.Position).Length()

	// For every object in the scene, check if the ray hits it.
	for _, prm := range s.Primitives {
		intersects, t, _, _ := prm.Geometry.RayCollision(r)
		if intersects && t < dtl{
			return true
		}
	}

	return false
}