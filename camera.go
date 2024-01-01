package gogl

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	Pos mgl32.Vec3

	Up      mgl32.Vec3
	Right   mgl32.Vec3
	Forward mgl32.Vec3

	WorldUp mgl32.Vec3

	Yaw   float32
	Pitch float32

	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func NewCamera(pos, worldUp mgl32.Vec3, yaw, pitch, speed, sensitivity float32) *Camera {
	cam := Camera{
		Pos:              pos,
		WorldUp:          worldUp,
		Yaw:              yaw,
		Pitch:            pitch,
		MovementSpeed:    speed,
		MouseSensitivity: sensitivity,
	}
	cam.updateVectors()

	return &cam
}

func (c *Camera) updateVectors() {
	forward := mgl32.Vec3{
		Cos32Deg(c.Yaw) * Cos32Deg(c.Pitch),
		Sin32Deg(c.Pitch),
		Sin32Deg(c.Yaw) * Cos32Deg(c.Pitch),
	}

	c.Forward = forward.Normalize()
	c.Right = forward.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Forward).Normalize()
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	center := c.Pos.Add(c.Forward)

	return mgl32.LookAt(
		c.Pos.X(), c.Pos.Y(), c.Pos.Z(),
		center.X(), center.Y(), center.Z(),
		c.Up.X(), c.Up.Y(), c.Up.Z(),
	)
}

func (c *Camera) UpdateCamera(dir MovementDirs, deltaTime, mouseDx, mouseDy float32) {
	magnitude := c.MovementSpeed * deltaTime

	//remove Z component and normalize
	forwardMovement := mgl32.Vec3{c.Forward.X(), 0, c.Forward.Z()}
	if forwardMovement.Len() > 0 {
		forwardMovement = forwardMovement.Normalize()
	}

	c.Pos = c.Pos.Add(forwardMovement.Mul(magnitude).Mul(float32(dir.Forward)))
	c.Pos = c.Pos.Add(c.Right.Mul(magnitude).Mul(float32(dir.Right)))
	c.Pos = c.Pos.Add(c.WorldUp.Mul(magnitude).Mul(float32(dir.Up)))

	mouseDx *= c.MouseSensitivity
	mouseDy *= c.MouseSensitivity

	c.Yaw += mouseDx
	if c.Yaw < 0 {
		c.Yaw = 360 - mgl32.Abs(c.Yaw)
	} else if c.Yaw >= 360 {
		c.Yaw -= 360
	}

	c.Pitch += mouseDy
	if c.Pitch >= 90 {
		c.Pitch = 89.9999
	} else if c.Pitch <= -90 {
		c.Pitch = -89.9999
	}

	c.updateVectors()
}

type MovementDirs struct {
	Forward int
	Right   int
	Up      int
}

func NewMoveDirs(f, b, r, l, u, d bool) MovementDirs {
	var fi, bi, ri, li, ui, di int
	if f {
		fi = 1
	}
	if b {
		bi = 1
	}
	if r {
		ri = 1
	}
	if l {
		li = 1
	}
	if u {
		ui = 1
	}
	if d {
		di = 1
	}

	moveDirs := MovementDirs{
		Forward: fi - bi,
		Right:   ri - li,
		Up:      ui - di,
	}
	return moveDirs
}
