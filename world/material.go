package world

import (
	"fmt"
)

// Material returns a initialized MaterialComponent
func Material() MaterialComponent {
	return MaterialComponent{}
}

// MaterialComponent of the thing
type MaterialComponent struct {
	DrawType uint
	// Program name (not being used yet)
	Program string
	// Dynamic for shaders?
	Props map[string]interface{}
}

// Material component
func (c *MaterialComponent) Material() *MaterialComponent { return c }

// Color color
type Color vec4

// RGBA produces html rgba color string
func (c Color) RGBA() string {
	return fmt.Sprintf("rgba(%d,%d,%d,%f)",
		byte(c[0]*255),
		byte(c[1]*255),
		byte(c[2]*255),
		c[3],
	)
}
