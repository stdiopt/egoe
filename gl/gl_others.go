// +build !js

package gl

// Types that others implements
type (
	Buffer       = Uint
	Shader       = Uint
	Program      = Uint
	Attrib       = Uint
	Framebuffer  = Uint
	Renderbuffer = Uint
	Texture      = Uint
	VertexArray  = Uint
	Uniform      = Uint
)

// Vars using those types
var (
	Buffer0 = Buffer(0)
)
