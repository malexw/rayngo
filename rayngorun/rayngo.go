package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"github.com/malexw/rayngo"
	"github.com/malexw/vmath"
)


func main() {
	width, height := 800, 480
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			ray := rayGen(x, y, width, height, 60)
			img.Set(x, height-y, background_color(ray))
		}
	}
	
	outFile, err := os.Create("test.png")
	if err != nil {
		fmt.Printf("File create failed.\n")
	}
	defer outFile.Close()
	png.Encode(outFile, img)
}

func rayGen(x int, y int, width int, height int, fov int) rayngo.Ray {
	// Calculate the eye position from the width, height, and field of view.
	// Assume {0,0,0} is the bottom left pixel of the screen
	eyeX := float64(width-1) / 2
	eyeY := float64(height-1) / 2
	eyeZ := eyeX / math.Tan(float64(fov)*math.Pi/180)
	eye := vmath.Vec3{eyeX, eyeY, eyeZ}

	// When returning the ray, make the origin some interesting location in world space.
	screenPos := vmath.Vec3{float64(x), float64(y), 0}
	return rayngo.Ray{vmath.Vec3{0,5,0}, screenPos.Sub(eye).Normalize()}
}

func background_color(ray rayngo.Ray) color.RGBA {
	// TODO Correct for vectors that point straight up
	
	dir := ray.Direction
	// Check for an intersection with the floor. If Y points downward, it must intersect eventually.
	if dir.Y < 0 {
		// Determine the colour of the floor at the intersection point.
		t := (0 - ray.Origin.Y) / dir.Y
		inter := ray.Origin.Add(dir.Scale(t))
		if math.Mod((math.Ceil(inter.X) + math.Ceil(inter.Z)), 2) == 0 {
			return color.RGBA{128, 0, 0, 255}
		} else {
			return color.RGBA{128, 128, 128, 255}
		}
	} else {
		// Otherwise, we're in the sky.
		angle := (math.Atan(dir.Y/math.Abs(dir.Z))/math.Pi) + 0.5
		return color.RGBA{0, 0, uint8(255 - (255 * angle)), 255}
	}
}