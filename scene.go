package rayngo

import (
	"image/color"
	"math"
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
	s.Shapes = append(s.Shapes, Shape{vmath.Vec3{3.0, 2.0, -11.0}, 1.0, vmath.Vec3{0, 128, 0}})
	s.Shapes = append(s.Shapes, Shape{vmath.Vec3{-4.0, 3.0, -9.0}, 2.0, vmath.Vec3{192, 192, 0}})

	return &s
}

func (s *Scene) BackgroundColor(r Ray) color.RGBA {
	// TODO Correct for vectors that point straight up
	
	dir := r.Direction
	// Check for an intersection with the floor. If Y points downward, it must intersect eventually.
	if dir.Y < 0 {
		// First find the intersection point with the floor.
		t := (0 - r.Origin.Y) / dir.Y
		inter := r.Origin.Add(dir.Scale(t))

		// Check to see if we're in shadow or not.
		sray := Ray{inter, s.LightSrc.Position.Sub(inter).Normalize()}
		inShadow := s.isShadowed(sray)
		brightness := 1.0
		if inShadow {
			brightness = 0.2
		}

		// Determine the colour of the floor at the intersection point.
		if math.Mod((math.Ceil(inter.X) + math.Ceil(inter.Z)), 2) == 0 {
			return color.RGBA{uint8(128*brightness), 0, 0, 255}
		} else {
			colorVal := uint8(128 * brightness)
			return color.RGBA{colorVal, colorVal, colorVal, 255}
		}
	} else {
		// Otherwise, we're in the sky.
		angle := (math.Atan(dir.Y/math.Abs(dir.Z))/math.Pi) + 0.5
		return color.RGBA{0, 0, uint8(255 - (255 * angle)), 255}
	}
}

func (s *Scene) isShadowed(r Ray) bool {
	// For every object in the scene, check if the ray hits it.
	for _, shp := range s.Shapes {
		intersects, _ := shp.RayCollision(r)
		if intersects {
			return true
		}
	}

	return false
}