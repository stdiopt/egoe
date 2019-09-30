package renderer

import (
	"errors"

	"github.com/stdiopt/egoe/gl"
)

func compileProgram(g gl.Context3, vertSrc, fragSrc string) (gl.Program, error) {
	var program gl.Program
	vertShader := g.CreateShader(gl.VERTEX_SHADER)
	g.ShaderSource(vertShader, vertShaderSrc)
	g.CompileShader(vertShader)
	if g.GetShaderi(vertShader, gl.COMPILE_STATUS) == 0 {
		return program, errors.New(g.GetShaderInfoLog(vertShader))
	}

	fragShader := g.CreateShader(gl.FRAGMENT_SHADER)
	g.ShaderSource(fragShader, fragShaderSrc)
	g.CompileShader(fragShader)

	if g.GetShaderi(fragShader, gl.COMPILE_STATUS) == 0 {
		return program, errors.New(g.GetShaderInfoLog(fragShader))
	}

	// Based on material
	program = g.CreateProgram()
	g.AttachShader(program, vertShader)
	g.AttachShader(program, fragShader)
	g.LinkProgram(program)
	if g.GetProgrami(program, gl.LINK_STATUS) == 0 {
		return program, errors.New(g.GetProgramInfoLog(program))
	}

	return program, nil
}

var vertShaderSrc = `#version 300 es
layout (location = 0) in vec3 a_position;
layout (location = 1) in vec4 a_color;
layout (location = 2) in mat4 a_transform;

uniform mat4 projection;

out vec4 colorV;

void main() {
	colorV = a_color; 
	gl_Position = projection * a_transform * vec4(a_position,1);
}
`

var fragShaderSrc = `#version 300 es
precision mediump float;

in vec4 colorV;
out vec4 color;

void main() {
	color = vec4(colorV.rgb * colorV.a, colorV.a); 
}
`
