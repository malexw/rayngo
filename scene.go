package rayngo

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
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

	// s.LightSrc = Light{vmath.Vec3{8.0, 7.0, 0.0}, ColorWhite(), 0.2, 0.7, 0.8}
	// s.Primitives = append(s.Primitives, *(&Primitive{}).AddShape(Sphere{vmath.Vec3{1.0, 2.5, -13.0}, 2.5}).SetMaterial(Material{Color{0, 0, 0.75, 1.0}, 50, false}))
	// s.Primitives = append(s.Primitives, *(&Primitive{}).AddShape(Sphere{vmath.Vec3{3.0, 1.0, -7.0}, 1.0}).SetMaterial(Material{Color{0, 0.5, 0, 1.0}, 20, false}))
	// s.Primitives = append(s.Primitives, *(&Primitive{}).AddShape(Sphere{vmath.Vec3{-4.0, 1.5, -9.0}, 1.5}).SetMaterial(Material{Color{0.75, 0, 0, 1.0}, 100, false}))
	// s.Primitives = append(s.Primitives, *(&Primitive{}).AddShape(Plane{vmath.Vec3{0.0, 1.0, 0.0}, 0.0}).SetMaterial(Material{Color{0.5, 0.5, 0.5, 1.0}, 10, true}))

	// pyrMat := *(NewMaterialFromMtl("res/materials/orange.mtl"))
	// pyramid := NewPrimitiveFromObj("res/meshes/pyramid.obj")
	// pyramid.SetMaterial(pyrMat)
	// pyramid.SetSqt(vmath.Vec3{8.0, 8.0, 8.0}, vmath.QuaternionFromAxisAngle(vmath.Vec3{0.0, 1.0, 0.0}, 50), vmath.Vec3{-10.0, 0.0, -21.0})
	// s.Primitives = append(s.Primitives, *pyramid)

	return &s
}

func NewSceneFromFile(path string) *Scene {
	fi, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var sd interface{}
	jerr := json.Unmarshal(fi, &sd)
	if jerr != nil {
		panic(jerr)
	}

	scene := Scene{
		Primitives: make([]Primitive, 0),
	}
	base := sd.(map[string]interface{})
	primitives := base["scene"].([]interface{})

	for _, prmdesc := range primitives {
		prmdict := prmdesc.(map[string]interface{})
		switch prmdict["type"] {
		case "sphere":
			radius := prmdict["radius"].(float64)
			material := prmdict["material"].(string)
			positionArray := prmdict["position"].([]interface{})
			position := vmath.Vec3{positionArray[0].(float64), positionArray[1].(float64), positionArray[2].(float64)}
			s := Sphere{position, radius}
			m := *(NewMaterialFromMtl(fmt.Sprintf("res/materials/%s", material)))
			scene.Primitives = append(scene.Primitives, *(&Primitive{}).AddShape(s).SetMaterial(m))
		case "plane":
			normalArray := prmdict["normal"].([]interface{})
			normal := vmath.Vec3{normalArray[0].(float64), normalArray[1].(float64), normalArray[2].(float64)}
			distance := prmdict["distance"].(float64)
			s := Plane{normal, distance}
			m := Material{Color{0.5, 0.5, 0.5, 1.0}, 10, true}
			scene.Primitives = append(scene.Primitives, *(&Primitive{}).AddShape(s).SetMaterial(m))
		case "mesh":
			meshPath := prmdict["mesh"].(string)
			mtlPath := prmdict["material"].(string)
			positionArray := prmdict["position"].([]interface{})
			position := vmath.Vec3{positionArray[0].(float64), positionArray[1].(float64), positionArray[2].(float64)}
			scaleArray := prmdict["scale"].([]interface{})
			scale := vmath.Vec3{scaleArray[0].(float64), scaleArray[1].(float64), scaleArray[2].(float64)}
			rot := prmdict["rotation"].(map[string]interface{})
			axisArray := rot["axis"].([]interface{})
			axis := vmath.Vec3{axisArray[0].(float64), axisArray[1].(float64), axisArray[2].(float64)}
			angle := rot["angle"].(float64)
			p := NewPrimitiveFromObj(fmt.Sprintf("res/meshes/%s", meshPath))
			m := *(NewMaterialFromMtl(fmt.Sprintf("res/materials/%s", mtlPath)))
			p.SetMaterial(m)
			p.SetSqt(scale, vmath.QuaternionFromAxisAngle(axis, angle), position)
			scene.Primitives = append(scene.Primitives, *p)
		}
	}

	scene.LightSrc = Light{vmath.Vec3{8.0, 7.0, 0.0}, ColorWhite(), 0.2, 0.7, 0.8}
	return &scene
}

func (s *Scene) BackgroundColor(r Ray) Color {
	// TODO Correct for vectors that point straight up

	dir := r.Direction
	angle := (math.Atan(dir.Y/math.Abs(dir.Z))/math.Pi) + 0.5
	return Color{0.55 - (0.55 * angle), 0.82 - (0.82 * angle), 1.0 - (1.0 * angle), 1.0}
}

func (s *Scene) IsShadowed(r Ray) bool {
	dtl := r.Origin.Sub(s.LightSrc.Position).Length()

	// For every object in the scene, check if the ray hits it.
	for _, prm := range s.Primitives {
		for shp := prm.Geometry.Front(); shp != nil; shp = shp.Next() {
			intersects, t, _, _ := shp.Value.(Shape).RayCollision(r)
			if intersects && t < dtl{
				return true
			}
		}
	}

	return false
}
