package gl

// Context3 opengles3/webgl2 funcs
// Not all
type Context3 interface {
	Context

	// glGetUniformBlockIndex retrieves the index of a uniform block within program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockIndex.xhtml
	GetUniformBlockIndex(p Program, name string) int

	// UniformBlockBinding assign a binding point to an active uniform block
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniformBlockBinding.xhtml
	UniformBlockBinding(p Program, index, bind int)

	// BindBufferBase bind a buffer object to an indexed buffer target
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBufferBase.xhtml
	BindBufferBase(target Enum, n uint32, b Buffer)

	// DrawArraysInstanced draw multiple instances of a range of elements
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
	DrawArraysInstanced(mode Enum, first int, count, primcount int)

	// VertexAttribDivisor  modify the rate at which generic vertex attributes
	// advance during instanced rendering
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArraysInstanced.xhtml
	VertexAttribDivisor(index Attrib, divisor int)

	// {lpf} Custom func
	// BufferDataX will type switch the interface and select the proper type
	BufferDataX(target Enum, d interface{}, usage Enum)
}
