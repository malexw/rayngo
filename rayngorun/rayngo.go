package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"github.com/malexw/pb"
	"github.com/malexw/rayngo"
	"github.com/malexw/vmath"
)


func main() {
	conf := rayngo.Config{480, 800, false, 64}
	pxCount := conf.Width * conf.Height

	dofAtten := 1.0 / float64(conf.DofRayCount)

	progressBar := pb.StartNew(pxCount)
	progressBar.ShowSpeed = true
	progressBar.MaxWidth = 100

	img := image.NewRGBA(image.Rect(0, 0, conf.Width, conf.Height))

	scene := rayngo.NewSceneFromFile("res/scene")

	for y := 0; y < conf.Height; y += 1 {
		for x := 0; x < conf.Width; x += 1 {
			accumulatedColor := rayngo.Color{0.0, 0.0, 0.0, 1.0}
			ray := rayGen(x, y, conf.Width, conf.Height, 40)
			// Final loop for depth of field. Do n samples per pixel and average the results
			if conf.DofEnabled {
				for n := 0; n < conf.DofRayCount; n += 1 {
					dofRay := dofRayGen(ray, 14)
					accumulatedColor = accumulatedColor.Add(collision(dofRay, scene).Attenuate(dofAtten))
				}
			} else {
				accumulatedColor = collision(ray, scene)
			}
			img.Set(x, conf.Height-y, accumulatedColor.ToImageColor())
			progressBar.Increment()
		}
	}

	progressBar.Finish()

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

func dofRayGen(srcRay rayngo.Ray, fDist float64) rayngo.Ray {
	fPoint := srcRay.Origin.Add(srcRay.Direction.Scale(fDist))

	xRand := (rand.Float64() - 0.5) * 0.2
	yRand := (rand.Float64() - 0.5) * 0.2
	pert := vmath.Vec3{xRand, yRand, 0.0}
	newSrc := srcRay.Origin.Add(pert)
	newDirection := fPoint.Sub(newSrc)

	return rayngo.Ray{newSrc, newDirection.Normalize()}
}

func collision(ray rayngo.Ray, scene *rayngo.Scene) rayngo.Color {
	c := rayngo.Color{0.0, 0.0, 0.0, 1.0}
	objHit := false
	nearest := math.MaxFloat64

	// For every object in the scene, check if the ray hits it. If it does, return the color of the object.
	for _, prm := range scene.Primitives {
		for shp := prm.Geometry.Front(); shp != nil; shp = shp.Next() {

			intersects, distance, location, norm := shp.Value.(rayngo.Shape).RayCollision(ray)
			if intersects && distance < nearest {
				nearest = distance
				// TODO Convert to texture coordinates, and use those to sample
				rawColor := prm.Mat.Sample(location.X, location.Z)
				// Ambient term
				ambColor := rawColor.Attenuate(scene.LightSrc.AmbientCoeff)

				// Check and see if the intersection point is in a shadow. If so, skip the diffuse and
				// specular terms.
				dirToLight := (scene.LightSrc.Position.Sub(location)).Normalize()
				// Scale by a small term in the direction of the light source to prevent intersections with self
				sray := rayngo.Ray{location.Add(dirToLight.Scale(0.0001)), dirToLight}
				inShadow := scene.IsShadowed(sray)

				if !inShadow {
					// Diffuse term
					lightness := math.Max(0, norm.Dot(dirToLight))
					diffColor := rawColor.Attenuate(scene.LightSrc.DiffuseCoeff).Attenuate(lightness)

					// Specular term
					reflected := norm.Scale(2*dirToLight.Dot(norm)).Sub(dirToLight)
					specular := math.Max(0, reflected.Dot(ray.Direction.Scale(-1)))

					// TODO Use LightSrc.Specular color once it exists instead of LightSrc.Diffuse
					specColor := scene.LightSrc.Diffuse.Attenuate(scene.LightSrc.SpecularCoeff)
					specColor = specColor.Attenuate(math.Pow(specular, float64(prm.Mat.Specularity)))

					c = ambColor.Add(diffColor).Add(specColor)
				} else {
					c = ambColor
				}
				objHit = true
			}
		}
	}

	// Otherwise, return the color of the background.
	if !objHit {
		c = scene.BackgroundColor(ray)
	}

	return c
}
