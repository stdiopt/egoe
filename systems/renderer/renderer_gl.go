// Package renderer mostly agnostic gl renderer with a couple of exceptions due
// to some limitations
package renderer

import (
	"github.com/stdiopt/egoe/ecs"
	"github.com/stdiopt/egoe/gl"
	"github.com/stdiopt/egoe/world"

	"github.com/go-gl/mathgl/mgl32"
)

type glRenderable interface {
	Transform() *world.TransformComponent
	Operation() *world.OperationComponent
	Material() *world.MaterialComponent
}

// Batch on add
type renderableInstance struct {
	Operation interface{}
	VAO       gl.VertexArray
	VBO       gl.Buffer
	TRO       gl.Buffer // transform buffer
	CLO       gl.Buffer // color buffers

	renderables []glRenderable
	PointsLen   int

	transBuf *F32TransferBuf
	colorBuf *F32TransferBuf
}

type glRenderer struct {
	// Camera projection
	projection mgl32.Mat4
	// Instances things

	instances []*renderableInstance

	g gl.Context3

	// Material
	program     gl.Program
	aPosition   gl.Attrib
	aTransform  gl.Attrib
	aColor      gl.Attrib
	uProjection gl.Uniform
}

// System initializes gl context and attatch the handlers on manager
func System(g gl.Context3, m *ecs.Manager) {
	g.ClearColor(1, 1, 1, 1)
	g.LineWidth(4)
	g.Enable(gl.CULL_FACE)
	g.CullFace(gl.BACK)

	rs := glRenderer{
		g:         g,
		instances: []*renderableInstance{},
	}
	rs.setupMaterial()
	m.ReflAuto(&rs)

}

func (rs *glRenderer) setupMaterial() {
	g := rs.g
	///////////////////////////////////////////////////////////////////////////
	// Generate material thing
	////////////////
	g.Enable(gl.BLEND)
	g.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	program, err := compileProgram(g, vertShaderSrc, fragShaderSrc)
	if err != nil {
		panic(err)
	}
	rs.program = program
	rs.aPosition = g.GetAttribLocation(program, "a_position")
	rs.aTransform = g.GetAttribLocation(program, "a_transform")
	rs.aColor = g.GetAttribLocation(program, "a_color")
	rs.uProjection = g.GetUniformLocation(program, "projection")
}

func (rs *glRenderer) WatchResize(evt world.ResizeEvent) {
	sz := mgl32.Vec2(evt)
	rs.g.Viewport(0, 0, int(sz[0]), int(sz[1]))
}
func (rs *glRenderer) WatchCamera(evt world.CameraEvent) {
	rs.projection = mgl32.Mat4(evt)
}

func (rs *glRenderer) HandleEntityAdd(evt ecs.EntitiesAddEvent) {
	g := rs.g
	entities := []ecs.Entity(evt)

	instances := []*renderableInstance{}

	var instanced *renderableInstance
	// Build the thing
	for _, e := range entities {
		r, ok := e.(glRenderable)
		if !ok {
			continue
		}
		operation := r.Operation()
		switch op := operation.OP.(type) {
		case *world.Poly:
			if instanced == nil || instanced.Operation != op {
				instanced = &renderableInstance{}
			}
			if instanced.Operation == nil { // Create a new
				instanced.VAO = g.CreateVertexArray()

				instanced.PointsLen = len(op.Points) / 3
				instanced.Operation = op

				g.BindVertexArray(instanced.VAO)

				// Create Vertex buffer And upload points
				instanced.VBO = g.CreateBuffer()
				g.BindBuffer(gl.ARRAY_BUFFER, instanced.VBO)
				g.BufferDataX(gl.ARRAY_BUFFER, op.Points, gl.STATIC_DRAW)
				g.EnableVertexAttribArray(rs.aPosition)
				g.VertexAttribPointer(rs.aPosition, 3, gl.FLOAT, false, 0, 0)

				// Create buffer and set the vertex pointers for the transforms
				instanced.TRO = g.CreateBuffer()
				g.BindBuffer(gl.ARRAY_BUFFER, instanced.TRO)

				// Hack to upload a mat4 into attr (4*4 vec4)
				vec4size := 4 * 4 // 4 floats in bytes size
				for i := uint32(0); i < 4; i++ {
					//log.Printf("Enable: %d sz: %d offset: %d", rs.aTransform+i, 4*vec4size, int(i)*vec4size)
					g.EnableVertexAttribArray(rs.aTransform + i)
					g.VertexAttribPointer(rs.aTransform+i, 4, gl.FLOAT, false, 4*vec4size, int(i)*vec4size)
					g.VertexAttribDivisor(rs.aTransform+i, 1)
				}

				// Create color buffer
				instanced.CLO = g.CreateBuffer()
				g.BindBuffer(gl.ARRAY_BUFFER, instanced.CLO)
				g.EnableVertexAttribArray(rs.aColor)
				g.VertexAttribPointer(rs.aColor, 4, gl.FLOAT, false, 0, 0)
				g.VertexAttribDivisor(rs.aColor, 1)

				instances = append(instances, instanced)
			}
			instanced.renderables = append(instanced.renderables, r)
		}
	}

	// Go through the instances and prepare those buffers
	for _, ins := range instances {
		sz := len(ins.renderables)
		ins.transBuf = NewF32TransferBuf(sz * 16)
		ins.colorBuf = NewF32TransferBuf(len(ins.renderables) * 4)
	}

	rs.instances = append(rs.instances, instances...)
}
func (rs *glRenderer) HandleUpdate(evt world.UpdateEvent) {
	g := rs.g
	g.Clear(gl.COLOR_BUFFER_BIT)

	for _, ins := range rs.instances {
		if len(ins.renderables) == 0 {
			continue
		}
		g.UseProgram(rs.program) // material
		g.UniformMatrix4fv(rs.uProjection, rs.projection[:])
		material := ins.renderables[0].Material()

		// Instancing upload all transforms into a float array
		for i, r := range ins.renderables {
			// Do the transformations
			transform := r.Transform()
			material := r.Material()
			m := transform.Mat4()
			ins.transBuf.WriteAt(m[:], i*16)
			//ins.transBuf.set(i*16, m[:])

			color := material.Props["color"].(world.Color)
			//ins.colorBuf.set(i*4, color[:])
			ins.colorBuf.WriteAt(color[:], i*4)
		}

		// Upload transform buffer
		g.BindBuffer(gl.ARRAY_BUFFER, ins.TRO)
		g.BufferDataX(gl.ARRAY_BUFFER, ins.transBuf.Get(), gl.DYNAMIC_DRAW)

		// Upload transform buffer
		g.BindBuffer(gl.ARRAY_BUFFER, ins.CLO)
		g.BufferDataX(gl.ARRAY_BUFFER, ins.colorBuf.Get(), gl.DYNAMIC_DRAW)

		g.BindVertexArray(ins.VAO)
		g.DrawArraysInstanced(uint32(material.DrawType), 0, ins.PointsLen, len(ins.renderables))

		// Renderable once
		/*for _, r := range ins.renderables {
			transform := r.Transform()
			material := r.Material()
			m := rs.projection.Mul4(transform.Mat4())
			g.UniformMatrix4fv(rs.uTransform, m[:])
			// Map GLMaterial with component
			color := material.Props["color"].(world.Color)
			g.Uniform4fv(rs.uColor, color[:])

			// Call once not several times (instancing)
			g.DrawArrays(uint32(material.DrawType), 0, ins.PointsLen)
		}*/
	}
}
