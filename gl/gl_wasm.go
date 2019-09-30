package gl

// Wasm implementation not everything is implemented as I've been implementing
// on demand

import (
	"reflect"
	"syscall/js"
	"unsafe"
)

type (
	Buffer       = js.Value
	Shader       = js.Value
	Program      = js.Value
	Attrib       = uint32
	Framebuffer  = js.Value
	Renderbuffer = js.Value
	Texture      = js.Value
	VertexArray  = js.Value
	Uniform      = js.Value
)

// Common to move things to gl
var (
	b8buf  = js.Global().Get("Uint8Array").New(8)
	b12buf = js.Global().Get("Uint8Array").New(12)
	b16buf = js.Global().Get("Uint8Array").New(16)
	b64buf = js.Global().Get("Uint8Array").New(64)

	// Float buffers of the byte ones
	f2buf  = js.Global().Get("Float32Array").New(b8buf.Get("buffer"))
	f3buf  = js.Global().Get("Float32Array").New(b12buf.Get("buffer"))
	f4buf  = js.Global().Get("Float32Array").New(b16buf.Get("buffer"))
	f16buf = js.Global().Get("Float32Array").New(b64buf.Get("buffer"))
)

// WebGL exposes the methods
type WebGL struct {
	webgl
}

// GetWebGL Return a js.Value webgl context
func GetWebGL(v js.Value) WebGL {
	return WebGL{webgl{v}}
}

type webgl struct {
	ctx js.Value
}

func (c webgl) ActiveTexture(texture Enum) {
	c.ctx.Call("activeTexture", texture)
}

func (c webgl) AttachShader(p Program, s Shader) {
	c.ctx.Call("attachShader", p, s)
}

func (c webgl) BindAttribLocation(p Program, a Attrib, name string) {
	c.ctx.Call("bindAttribLocation", p, a, name)
}

func (c webgl) BindBuffer(target Enum, b Buffer) {
	c.ctx.Call("bindBuffer", target, b)
}

func (c webgl) BindFramebuffer(target Enum, fb Framebuffer) {
	c.ctx.Call("bindFramebuffer", target, fb)
}

func (c webgl) BindRenderbuffer(target Enum, rb Renderbuffer) {
	c.ctx.Call("bindRenderbuffer", target, rb)
}

func (c webgl) BindTexture(target Enum, t Texture) {
	c.ctx.Call("bindTexture", target, t)
}

func (c webgl) BindVertexArray(rb VertexArray) {
	c.ctx.Call("bindVertexArray", rb)
}

func (c webgl) BlendColor(red, green, blue, alpha float32) {
	c.ctx.Call("blendColor", red, green, blue, alpha)
}

func (c webgl) BlendEquation(mode Enum) {
	c.ctx.Call("blendEquation", mode)
}

func (c webgl) BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	c.ctx.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

func (c webgl) BlendFunc(sfactor, dfactor Enum) {
	c.ctx.Call("blendFunc", sfactor, dfactor)
}

func (c webgl) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	panic("not implemented")
}

func (c webgl) BufferData(target Enum, src []byte, usage Enum) {
	d := js.Global().Get("Uint8Array").New()
	js.CopyBytesToJS(d, src)
	c.ctx.Call("bufferData", target, d, usage)
}

func (c webgl) BufferInit(target Enum, size int, usage Enum) {
	panic("not implemented")
}

func (c webgl) BufferSubData(target Enum, offset int, data []byte) {
	panic("not implemented")
}

func (c webgl) CheckFramebufferStatus(target Enum) Enum {
	panic("not implemented")
}

func (c webgl) Clear(mask Enum) {
	c.ctx.Call("clear", mask)
}

func (c webgl) ClearColor(red, green, blue, alpha float32) {
	c.ctx.Call("clearColor", red, green, blue, alpha)
}

func (c webgl) ClearDepthf(d float32) {
	panic("not implemented")
}

func (c webgl) ClearStencil(s int) {
	panic("not implemented")
}

func (c webgl) ColorMask(red, green, blue, alpha bool) {
	c.ctx.Call("colorMask", red, green, blue, alpha)
}

func (c webgl) CompileShader(s Shader) {
	c.ctx.Call("compileShader", s)
}

func (c webgl) CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	panic("not implemented")
}

func (c webgl) CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	panic("not implemented")
}

func (c webgl) CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	panic("not implemented")
}

func (c webgl) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	panic("not implemented")
}

func (c webgl) CreateBuffer() Buffer {
	return Buffer(c.ctx.Call("createBuffer"))
}

func (c webgl) CreateFramebuffer() Framebuffer {
	panic("not implemented")
}

func (c webgl) CreateProgram() Program {
	return Program(c.ctx.Call("createProgram"))
}

func (c webgl) CreateRenderbuffer() Renderbuffer {
	panic("not implemented")
}

func (c webgl) CreateShader(ty Enum) Shader {
	return Shader(c.ctx.Call("createShader", ty))
}

func (c webgl) CreateTexture() Texture {
	panic("not implemented")
}

func (c webgl) CreateVertexArray() VertexArray {
	return VertexArray(c.ctx.Call("createVertexArray"))
}

func (c webgl) CullFace(mode Enum) {
	c.ctx.Call("cullFace", mode)
}

func (c webgl) DeleteBuffer(v Buffer) {
	c.ctx.Call("deleteBuffer", v)
}

func (c webgl) DeleteFramebuffer(v Framebuffer) {
	c.ctx.Call("deleteFramebuffer", v)
}

func (c webgl) DeleteProgram(p Program) {
	c.ctx.Call("deleteProgram", p)
}

func (c webgl) DeleteRenderbuffer(v Renderbuffer) {
	panic("not implemented")
}

func (c webgl) DeleteShader(s Shader) {
	panic("not implemented")
}

func (c webgl) DeleteTexture(v Texture) {
	panic("not implemented")
}

func (c webgl) DeleteVertexArray(v VertexArray) {
	panic("not implemented")
}

func (c webgl) DepthFunc(fn Enum) {
	panic("not implemented")
}

func (c webgl) DepthMask(flag bool) {
	panic("not implemented")
}

func (c webgl) DepthRangef(n, f float32) {
	panic("not implemented")
}

func (c webgl) DetachShader(p Program, s Shader) {
	panic("not implemented")
}

func (c webgl) Disable(cap Enum) {
	c.ctx.Call("disable", cap)
}

func (c webgl) DisableVertexAttribArray(a Attrib) {
	c.ctx.Call("disableVertexAttribArray", a)
}

func (c webgl) DrawArrays(mode Enum, first, count int) {
	c.ctx.Call("drawArrays", mode, first, count)
}

func (c webgl) DrawElements(mode Enum, count int, ty Enum, offset int) {
	c.ctx.Call("drawElements", mode, count, ty, offset)
}

func (c webgl) Enable(cp Enum) {
	c.ctx.Call("enable", cp)
}

func (c webgl) EnableVertexAttribArray(a Attrib) {
	c.ctx.Call("enableVertexAttribArray", a)
}

func (c webgl) Finish() {
	panic("not implemented")
}

func (c webgl) Flush() {
	panic("not implemented")
}

func (c webgl) FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	panic("not implemented")
}

func (c webgl) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	panic("not implemented")
}

func (c webgl) FrontFace(mode Enum) {
	c.ctx.Call("frontFace", mode)
}

func (c webgl) GenerateMipmap(target Enum) {
	panic("not implemented")
}

func (c webgl) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	res := c.ctx.Call("getActiveAttrib", p, index)
	return res.Get("name").String(), res.Get("size").Int(), Enum(res.Get("type").Int())
}

func (c webgl) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	panic("not implemented")
}

func (c webgl) GetAttachedShaders(p Program) []Shader {
	panic("not implemented")
}

func (c webgl) GetAttribLocation(p Program, name string) Attrib {
	return Attrib(c.ctx.Call("getAttribLocation", p, name).Int())
}

func (c webgl) GetBooleanv(dst []bool, pname Enum) {
	panic("not implemented")
}

func (c webgl) GetFloatv(dst []float32, pname Enum) {
	panic("not implemented")
}

func (c webgl) GetIntegerv(dst []int32, pname Enum) {
	panic("not implemented")
}

func (c webgl) GetInteger(pname Enum) int {
	panic("not implemented")
}

func (c webgl) GetBufferParameteri(target, value Enum) int {
	panic("not implemented")
}

func (c webgl) GetError() Enum {
	panic("not implemented")
}

func (c webgl) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	panic("not implemented")
}

func (c webgl) GetProgrami(p Program, pname Enum) int {
	r := c.ctx.Call("getProgramParameter", p, pname)

	switch r.Type() {
	case js.TypeBoolean:
		if r.Bool() {
			return 1
		}
		return 0
	case js.TypeNumber:
		return r.Int()
	default:
		panic("unknown type")
	}
}

func (c webgl) GetProgramInfoLog(p Program) string {
	return c.ctx.Call("getProgramInfoLog").String()
}

func (c webgl) GetRenderbufferParameteri(target, pname Enum) int {
	panic("not implemented")
}

func (c webgl) GetShaderi(s Shader, pname Enum) int {
	if c.ctx.Call("getShaderParameter", s, pname).Bool() {
		return 1
	}
	return 0
}

func (c webgl) GetShaderInfoLog(s Shader) string {
	return c.ctx.Call("getShaderInfoLog", s).String()
}

func (c webgl) GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	panic("not implemented")

}

func (c webgl) GetShaderSource(s Shader) string {
	panic("not implemented")
}

func (c webgl) GetString(pname Enum) string {
	panic("not implemented")
}

func (c webgl) GetTexParameterfv(dst []float32, target, pname Enum) {
	panic("not implemented")
}

func (c webgl) GetTexParameteriv(dst []int32, target, pname Enum) {
	panic("not implemented")
}

func (c webgl) GetUniformfv(dst []float32, src Uniform, p Program) {
	c.ctx.Call("getUniformfv", dst, src, p)
}

func (c webgl) GetUniformiv(dst []int32, src Uniform, p Program) {
	panic("not implemented")
}

func (c webgl) GetUniformLocation(p Program, name string) Uniform {
	return Uniform(c.ctx.Call("getUniformLocation", p, name))
}

func (c webgl) GetVertexAttribf(src Attrib, pname Enum) float32 {
	panic("not implemented")
}

func (c webgl) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	panic("not implemented")
}

func (c webgl) GetVertexAttribi(src Attrib, pname Enum) int32 {
	panic("not implemented")
}

func (c webgl) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	panic("not implemented")
}

func (c webgl) Hint(target, mode Enum) {
	panic("not implemented")
}

func (c webgl) IsBuffer(b Buffer) bool {
	panic("not implemented")
}

func (c webgl) IsEnabled(cap Enum) bool {
	panic("not implemented")

}

func (c webgl) IsFramebuffer(fb Framebuffer) bool {
	panic("not implemented")
}

func (c webgl) IsProgram(p Program) bool {
	panic("not implemented")
}

func (c webgl) IsRenderbuffer(rb Renderbuffer) bool {
	panic("not implemented")
}

func (c webgl) IsShader(s Shader) bool {
	panic("not implemented")
}

func (c webgl) IsTexture(t Texture) bool {
	panic("not implemented")
}

func (c webgl) LineWidth(width float32) {
	c.ctx.Call("lineWidth", width)
}

func (c webgl) LinkProgram(p Program) {
	c.ctx.Call("linkProgram", p)
}

func (c webgl) PixelStorei(pname Enum, param int32) {
	panic("not implemented")
}

func (c webgl) PolygonOffset(factor, units float32) {
	panic("not implemented")
}

func (c webgl) ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	panic("not implemented")
}

func (c webgl) ReleaseShaderCompiler() {
	panic("not implemented")
}

func (c webgl) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	panic("not implemented")
}

func (c webgl) SampleCoverage(value float32, invert bool) {
	panic("not implemented")
}

func (c webgl) Scissor(x, y, width, height int32) {
	panic("not implemented")
}

func (c webgl) ShaderSource(s Shader, src string) {
	c.ctx.Call("shaderSource", s, src)
}

func (c webgl) StencilFunc(fn Enum, ref int, mask uint32) {
	panic("not implemented")
}

func (c webgl) StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	panic("not implemented")
}

func (c webgl) StencilMask(mask uint32) {
	panic("not implemented")
}

func (c webgl) StencilMaskSeparate(face Enum, mask uint32) {
	panic("not implemented")
}

func (c webgl) StencilOp(fail, zfail, zpass Enum) {
	panic("not implemented")
}

func (c webgl) StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	panic("not implemented")
}

func (c webgl) TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	panic("not implemented")

}

func (c webgl) TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	panic("not implemented")

}

func (c webgl) TexParameterf(target, pname Enum, param float32) {
	panic("not implemented")
}
func (c webgl) TexParameterfv(target, pname Enum, params []float32) {
	panic("not implemented")
}
func (c webgl) TexParameteri(target, pname Enum, param int) {
	panic("not implemented")
}
func (c webgl) TexParameteriv(target, pname Enum, params []int32) {
	panic("not implemented")
}
func (c webgl) Uniform1f(dst Uniform, v float32) {
	c.ctx.Call("uniform1f", dst, v)
}
func (c webgl) Uniform1fv(dst Uniform, src []float32) {
	panic("not implemented")
}
func (c webgl) Uniform1i(dst Uniform, v int) {
	panic("not implemented")
}
func (c webgl) Uniform1iv(dst Uniform, src []int32) {
	panic("not implemented")
}
func (c webgl) Uniform2f(dst Uniform, v0, v1 float32) {
	c.ctx.Call("uniform2f", dst, v0, v1)
}
func (c webgl) Uniform2fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b8buf, F32Bytes(src...))
	c.ctx.Call("uniform2fv", dst, f2buf)
}
func (c webgl) Uniform2i(dst Uniform, v0, v1 int) {
	panic("not implemented")
}
func (c webgl) Uniform2iv(dst Uniform, src []int32) {
	panic("not implemented")
}
func (c webgl) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	c.ctx.Call("uniform3f", dst, v0, v1, v2)
}
func (c webgl) Uniform3fv(dst Uniform, src []float32) {
	panic("not implemented")
}
func (c webgl) Uniform3i(dst Uniform, v0, v1, v2 int32) {
	panic("not implemented")
}
func (c webgl) Uniform3iv(dst Uniform, src []int32) {
	panic("not implemented")
}
func (c webgl) Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	c.ctx.Call("uniform4f", dst, v0, v1, v2, v3)
}
func (c webgl) Uniform4fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b16buf, F32Bytes(src...))
	c.ctx.Call("uniform4fv", dst, f4buf)
}
func (c webgl) Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	panic("not implemented")
}
func (c webgl) Uniform4iv(dst Uniform, src []int32) {
	panic("not implemented")
}
func (c webgl) UniformMatrix2fv(dst Uniform, src []float32) {
	panic("not implemented")
}
func (c webgl) UniformMatrix3fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b64buf, F32Bytes(src...))
	c.ctx.Call("uniformMatrix3fv", dst, false, f16buf)
}
func (c webgl) UniformMatrix4fv(dst Uniform, src []float32) {
	js.CopyBytesToJS(b64buf, F32Bytes(src...))
	c.ctx.Call("uniformMatrix4fv", dst, false, f16buf)
}
func (c webgl) UseProgram(p Program) {
	c.ctx.Call("useProgram", p)
}
func (c webgl) ValidateProgram(p Program) {
	panic("not implemented")
}
func (c webgl) VertexAttrib1f(dst Attrib, x float32) {
	panic("not implemented")
}
func (c webgl) VertexAttrib1fv(dst Attrib, src []float32) {
	panic("not implemented")
}
func (c webgl) VertexAttrib2f(dst Attrib, x, y float32) {
	panic("not implemented")
}
func (c webgl) VertexAttrib2fv(dst Attrib, src []float32) {
	panic("not implemented")
}
func (c webgl) VertexAttrib3f(dst Attrib, x, y, z float32) {
	panic("not implemented")
}
func (c webgl) VertexAttrib3fv(dst Attrib, src []float32) {
	panic("not implemented")
}
func (c webgl) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	panic("not implemented")
}
func (c webgl) VertexAttrib4fv(dst Attrib, src []float32) {
	panic("not implemented")
}
func (c webgl) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	c.ctx.Call("vertexAttribPointer", dst, size, ty, normalized, stride, offset)
}
func (c webgl) Viewport(x, y, width, height int) {
	c.ctx.Call("viewport", x, y, width, height)
}

///////////////////////////////////////////////////////////////////////////////
// webgl 2 + extra funcs
///////////////////////////////////////////////////////////////////////////////
func (c webgl) BufferDataX(target Enum, d interface{}, usage Enum) {
	c.ctx.Call("bufferData", target, conv(d), usage)
}
func (c webgl) GetUniformBlockIndex(p Program, name string) int {
	return c.ctx.Call("getUniformBlockIndex", p, name).Int()
}
func (c webgl) UniformBlockBinding(p Program, index, bind int) {
	c.ctx.Call("uniformBlockBinding", p, index, bind)
}

func (c webgl) BindBufferBase(target Enum, n uint32, b Buffer) {
	c.ctx.Call("bindBufferBase", target, n, b)
}

func (c webgl) DrawArraysInstanced(mode Enum, first int, count, primcount int) {
	c.ctx.Call("drawArraysInstanced", mode, first, count, primcount)
}

func (c webgl) VertexAttribDivisor(index Attrib, divisor int) {
	c.ctx.Call("vertexAttribDivisor", index, divisor)
}

// F32Bytes unsafe cast list of floats to byte
func F32Bytes(values ...float32) []byte {
	// size in bytes
	f32size := 4
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&values))
	header.Len *= f32size
	header.Cap *= f32size

	// Convert slice header to []byte
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}

// conv will convert a slice to a typedarray
//  []float32 -> Float32Array
//  []float64 -> Float32Array (for webgl purposes)
func conv(data interface{}) js.Value {
	switch data := data.(type) {
	case js.Value:
		return data
	case []float32:
		d := newFloat32Array(len(data))
		for i, v := range data {
			d.SetIndex(i, v)
		}
		return d
	case []float64:
		d := newFloat32Array(len(data))
		for i, v := range data {
			d.SetIndex(i, float32(v))
		}
		return d
	default:
		panic("Unimplemented type")
	}

	return js.Undefined()

}

func newFloat32Array(sz int) js.Value {
	return js.Global().Get("Float32Array").New(sz)
}
