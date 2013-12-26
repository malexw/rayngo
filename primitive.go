package rayngo

import (
	"container/list"
	"github.com/malexw/vmath"
)

type Primitive struct {
	Geometry list.List
	ModelGeometry list.List
	Mat Material
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