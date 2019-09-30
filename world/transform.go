package world

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Affine for each
type Affine struct {
	Translation vec3
	Rotation    vec3
	Scale       vec3
}

// Mat4 returns the mat4 from the Affine
// S * R * T
func (c Affine) Mat4() mat4 {
	m := mgl32.Scale3D(c.Scale[0], c.Scale[1], c.Scale[2])

	if c.Rotation[0] != 0 {
		m = m.Mul4(mgl32.HomogRotate3DX(c.Rotation[0]))
	}
	if c.Rotation[1] != 0 {
		m = m.Mul4(mgl32.HomogRotate3DY(c.Rotation[1]))
	}
	if c.Rotation[2] != 0 {
		m = m.Mul4(mgl32.HomogRotate3DZ(c.Rotation[2]))
	}

	return m.Mul4(mgl32.Translate3D(c.Translation[0], c.Translation[1], c.Translation[2]))
}

// TransformComponent Thing
type TransformComponent struct {
	Local Affine
	World Affine
}

// Mat4 returns the mat4 from the transformations
// Multiply world * local
func (c *TransformComponent) Mat4() mat4 {
	return c.World.Mat4().
		Mul4(c.Local.Mat4())
}

// Transform returns an initialized transform component
func Transform() TransformComponent {
	return TransformComponent{
		Local: Affine{Scale: vec3{1, 1, 1}},
		World: Affine{Scale: vec3{1, 1, 1}},
	}
}

// Transform component
func (c *TransformComponent) Transform() *TransformComponent { return c }
