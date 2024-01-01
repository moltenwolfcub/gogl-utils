package helpers

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type BufferID uint32

// used for generating and binding general buffers
// e.g. VertexBufferObject or NormalArrayObject
func GenBindBuffer(target uint32) BufferID {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	return BufferID(buffer)
}

// used for generating and binding vertex buffers
func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func BindVertexArray(id BufferID) {
	gl.BindVertexArray(uint32(id))
}

// used for initialising the target with the go array at data
// a more go-esque wrapper for the gl.BufferData function
func BufferData[T any](target uint32, data []T, usage uint32) {
	var v T
	dataTypeSize := unsafe.Sizeof(v)

	gl.BufferData(target, len(data)*int(dataTypeSize), gl.Ptr(data), usage)
}

// change to interface
type BufferLayout[T any] struct {
	data     []T
	segments []int32
}

func NewBufferLayout[T any](segments []int32, data []T) BufferLayout[T] {
	b := BufferLayout[T]{
		data:     data,
		segments: segments,
	}
	return b
}

func (b BufferLayout[T]) dataSize() int32 {
	var v T
	dataTypeSize := unsafe.Sizeof(v)
	return int32(dataTypeSize)
}

func (b BufferLayout[T]) calcStride() int32 {
	var objLen int32
	for _, i := range b.segments {
		objLen += i
	}
	return objLen * b.dataSize()
}
func (b BufferLayout[T]) offset(index int) uintptr {
	var offset int32
	for i := index - 1; i >= 0; i-- {
		offset += b.segments[i]
	}
	return uintptr(offset * b.dataSize())
}
func (b BufferLayout[T]) getXType() uint32 {
	var v T
	i := reflect.TypeOf(v).Kind()
	switch i {
	case reflect.Float32:
		return gl.FLOAT
	default:
		panic(fmt.Errorf("unsupported type (%v) for generating xtype", i))
		//I've only implemented the one's I've had to use
	}
}

type BufferLoader struct {
	layoutIndex uint32
}

func NewBufferLoader() *BufferLoader {
	b := BufferLoader{}
	return &b
}

func (b *BufferLoader) BuildFloatBuffer(id BufferID, layout BufferLayout[float32]) {
	BufferData(gl.ARRAY_BUFFER, layout.data, gl.STATIC_DRAW)

	BindVertexArray(id)

	for i, s := range layout.segments {
		index := b.layoutIndex + uint32(i)
		gl.VertexAttribPointerWithOffset(index, s, layout.getXType(), false, layout.calcStride(), layout.offset(i))
		gl.EnableVertexAttribArray(index)
	}

	b.layoutIndex += uint32(len(layout.segments))
}

// VAO := helpers.GenBindVertexArray()
// helpers.BufferData(gl.ARRAY_BUFFER, verticies, gl.STATIC_DRAW)

// helpers.BindVertexArray(VAO)
// gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
// gl.EnableVertexAttribArray(0)
// gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, uintptr(3*4))
// gl.EnableVertexAttribArray(1)

// NAO := helpers.GenBindBuffer(gl.ARRAY_BUFFER)
// helpers.BufferData(gl.ARRAY_BUFFER, normals, gl.STATIC_DRAW)

// helpers.BindVertexArray(NAO)
// gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, nil)
// gl.EnableVertexAttribArray(2)
