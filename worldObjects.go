package gogl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	Type         string
	Verticies    []float32 //in XYZ UV
	VertexStride int       // 5 if using XYZ UV
	normals      []float32
	bufferLoader *BufferLoader
	vao          BufferID
	nao          BufferID
}

func (o *Object) FillBuffers() {
	o.bufferLoader = NewBufferLoader()
	o.vao = GenBindVertexArray()
	o.nao = GenBindBuffer(gl.ARRAY_BUFFER)

	GenBindBuffer(gl.ARRAY_BUFFER) //VBO

	BindVertexArray(o.vao)
	o.bufferLoader.BuildFloatBuffer(o.vao, NewBufferLayout([]int32{3, 2}, o.Verticies))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(o.nao))
	o.bufferLoader.BuildFloatBuffer(o.nao, NewBufferLayout([]int32{3}, o.normals))
}

func (o *Object) CalcNormals(triangleCount int) {
	vertexCount := triangleCount * 3 //3 bc we are working in 3d space so XYZ

	o.normals = make([]float32, vertexCount*3)
	for tri := 0; tri < triangleCount; tri++ {
		index := tri * o.VertexStride * 3
		p1 := mgl32.Vec3{o.Verticies[index], o.Verticies[index+1], o.Verticies[index+2]}
		index += o.VertexStride
		p2 := mgl32.Vec3{o.Verticies[index], o.Verticies[index+1], o.Verticies[index+2]}
		index += o.VertexStride
		p3 := mgl32.Vec3{o.Verticies[index], o.Verticies[index+1], o.Verticies[index+2]}

		normal := TriangleNormal(p1, p2, p3)
		o.normals[tri*9+0] = normal.X()
		o.normals[tri*9+1] = normal.Y()
		o.normals[tri*9+2] = normal.Z()

		o.normals[tri*9+3] = normal.X()
		o.normals[tri*9+4] = normal.Y()
		o.normals[tri*9+5] = normal.Z()

		o.normals[tri*9+6] = normal.X()
		o.normals[tri*9+7] = normal.Y()
		o.normals[tri*9+8] = normal.Z()
	}
}

func (o Object) Draw(shader Shader, drawMatrix mgl32.Mat4) {
	BindVertexArray(o.vao)

	shader.SetMatrix4("model", drawMatrix)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.Verticies)/o.VertexStride))
}

func (o Object) DrawMultiple(shader Shader, num int, drawMatrix func(int) mgl32.Mat4) {
	BindVertexArray(o.vao)

	for i := 0; i < num; i++ {
		shader.SetMatrix4("model", drawMatrix(i))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.Verticies)/o.VertexStride))
	}
}

func Cube(size float32) Object {
	o := Object{
		Type: "cube",
	}
	o.Verticies = []float32{
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,
		-size / 2, size / 2, -size / 2, 0.0, 1.0,

		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, size / 2, size / 2, 1.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 1.0,
		-size / 2, size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,

		-size / 2, size / 2, size / 2, 1.0, 0.0,
		-size / 2, size / 2, -size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		-size / 2, size / 2, size / 2, 1.0, 0.0,

		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 0.0,

		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,

		-size / 2, size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		-size / 2, size / 2, -size / 2, 0.0, 1.0,
		-size / 2, size / 2, size / 2, 0.0, 0.0,
	}
	o.VertexStride = 5

	o.CalcNormals(12)
	o.FillBuffers()

	return o
}

func Pentahedron(size float32) Object {
	o := Object{
		Type: "pentahedron",
	}
	o.Verticies = []float32{
		size / 2, -size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 0.0,
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		-size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
	}
	o.VertexStride = 5

	o.CalcNormals(6)
	o.FillBuffers()

	return o
}
