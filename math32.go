package helpers

/*
Type conversions make the main code look messy
so some wrapper functions to make it more readible
*/

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func Sin32(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Cos32(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sin32Deg(x float32) float32 {
	return Sin32(mgl32.DegToRad(x))
}

func Cos32Deg(x float32) float32 {
	return Cos32(mgl32.DegToRad(x))
}

func Mod32(a, b float32) float32 {
	return float32(math.Mod(float64(a), float64(b)))
}
