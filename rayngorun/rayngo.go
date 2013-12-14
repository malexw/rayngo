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

	//sphere := rayngo.Shape{vmath.Vec3{3.0, 5.0, -10.0}, 1.0, color.RGBA{128, 128, 0, 255}}
	scene := rayngo.NewScene()

	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			ray := rayGen(x, y, width, height, 40)
			img.Set(x, height-y, collision(ray, scene))
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
	return rayngo.Ray{vmath.Vec3{0,3,7}, screenPos.Sub(eye).Normalize()}
}

func collision(ray rayngo.Ray, scene *rayngo.Scene) color.RGBA {
	c := color.RGBA{0,0,0,255}
	objHit := false

	// For every object in the scene, check if the ray hits it. If it does, return the color of the object.
	for _, prm := range scene.Primitives {
		intersects, location := prm.Geometry.RayCollision(ray)
		if intersects {
			// Ambient term
			ambColor := prm.Mat.Diffuse.Attenuate(scene.LightSrc.AmbientCoeff)

			// Diffuse term
			norm := (location.Sub(prm.Geometry.Position)).Normalize()
			dirToLight := (scene.LightSrc.Position.Sub(location)).Normalize()
			lightness := math.Max(0, norm.Dot(dirToLight))
			diffColor := prm.Mat.Diffuse.Attenuate(scene.LightSrc.DiffuseCoeff).Attenuate(lightness)

			// Specular term
			reflected := norm.Scale(2*dirToLight.Dot(norm)).Sub(dirToLight)
			specular := math.Max(0, reflected.Dot(ray.Direction.Scale(-1)))

			// TODO Use LightSrc.Specular color once it exists instead of LightSrc.Diffuse
			specColor := scene.LightSrc.Diffuse.Attenuate(scene.LightSrc.SpecularCoeff)
			specColor = specColor.Attenuate(math.Pow(specular, float64(prm.Mat.Specularity)))

			c = ambColor.Add(diffColor).Add(specColor).ToImageColor()
			objHit = true
		}
	}

	// Otherwise, return the color of the background.
	if !objHit {
		c = scene.BackgroundColor(ray)
	}

	return c
}

func vec3ToColor(v vmath.Vec3) color.RGBA {
	return color.RGBA{uint8(math.Min(v.X, 1.0)*255), uint8(math.Min(v.Y, 1.0)*255), uint8(math.Min(v.Z, 1.0)*255), 255}
}