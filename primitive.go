package rayngo

import (
	"bufio"
	"container/list"
	"os"
	"strconv"
	"strings"
	"github.com/malexw/vmath"
)

type Primitive struct {
	Geometry list.List
	ModelGeometry list.List
	Mat Material
}

func NewPrimitiveFromObj(path string) *Primitive {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var prm Primitive
	var vertices []vmath.Vec3

	lineScanner := bufio.NewScanner(fi)

	for lineScanner.Scan() {
		line := lineScanner.Text()

		var tokens []string
		wordScanner := bufio.NewScanner(strings.NewReader(line))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			tokens = append(tokens, wordScanner.Text())
		}

		switch tokens[0] {
		case "v":
			f1, _ := strconv.ParseFloat(tokens[1], 64)
			f2, _ := strconv.ParseFloat(tokens[2], 64)
			f3, _ := strconv.ParseFloat(tokens[3], 64)
			vertices = append(vertices, vmath.Vec3{f1, f2, f3})
		// case "vt":
		// 	do nothing. Don't support texture mapping yet.
		// case "vn":
		// 	do nothing. Calculated automatically.
		// 	normals = append(normals, vmath.Vec3{strconv.ParseFloat(tokens[1]), strconv.ParseFloat(tokens[2]), strconv.ParseFloat(tokens[3])})
		case "f":
			idx1, _ := strconv.ParseInt(strings.Split(tokens[1], "/")[0], 10, 64)
			idx2, _ := strconv.ParseInt(strings.Split(tokens[2], "/")[0], 10, 64)
			idx3, _ := strconv.ParseInt(strings.Split(tokens[3], "/")[0], 10, 64)
			tri := Triangle{[3]vmath.Vec3{vertices[idx1 - 1], vertices[idx2 - 1], vertices[idx3 - 1]}}
			prm.AddShape(tri)
		// case "usemtl":
		// 	do nothing. Parsing mtl files isn't supported yet.
		}
	}

	return &prm
}

func (p *Primitive) AddShape(s Shape) *Primitive {
	p.ModelGeometry.PushBack(s)
	p.Geometry.PushBack(s)
	return p
}

func (p *Primitive) SetMaterial(m Material) *Primitive {
	p.Mat = m
	return p
}

func (p *Primitive) SetSqt(scale vmath.Vec3, rot vmath.Quaternion, trans vmath.Vec3) *Primitive {
	p.Geometry = list.List{}
	xform := vmath.Matrix4FromSqt(scale, rot, trans)

	for shape := p.ModelGeometry.Front(); shape != nil; shape = shape.Next() {
		p.Geometry.PushBack(shape.Value.(Shape).Transform(xform))
	}

	return p
}
