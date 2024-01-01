package gogl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func TriangleNormal(p1, p2, p3 mgl32.Vec3) mgl32.Vec3 {
	U := p2.Sub(p1)
	V := p3.Sub(p1)

	return U.Cross(V).Normalize()
}
