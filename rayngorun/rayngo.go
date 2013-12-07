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
			img.Set(x, y, background_color(ray.Direction))
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
	// Assume {0,0,0} is the bottom left pixel of the screen, 
	eyeX := float64(width-1) / 2
	eyeY := float64(height-1) / 2
	eyeZ := eyeX / math.Tan(float64(fov)*math.Pi/180)
	eye := vmath.Vec3{eyeX, eyeY, eyeZ}

	// When returning the ray, we want the eye position to be the origin.
	screenPos := vmath.Vec3{float64(x), float64(y), 0}
	return rayngo.Ray{vmath.Vec3{0,0,0}, screenPos.Sub(eye).Normalize()}
}

func background_color(dir vmath.Vec3) color.RGBA {
	// TODO Correct for vectors that point straight up
	angle := (math.Atan(float64(dir.Y)/math.Abs(float64(dir.Z)))/math.Pi) + 0.5
	return color.RGBA{0, 0, uint8(255 * angle), 255}
}