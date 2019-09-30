// +build js,wasm

package renderer

import (
	"reflect"
	"syscall/js"
	"unsafe"
)

const f32size = 4

//F32TransferBuf holds a buffer and a js reference
type F32TransferBuf struct {
	// Buffer is the local go buffer that we write
	buffer  []float32
	jsBytes js.Value
}

//NewF32TransferBuf creates a go2js transfer buf
func NewF32TransferBuf(sz int) *F32TransferBuf {
	buffer := make([]float32, sz)
	jsBytes := js.Global().Get("Uint8Array").New(sz * f32size)
	return &F32TransferBuf{
		buffer:  buffer,
		jsBytes: jsBytes,
	}
}

// WriteAt write floats at offset in buffer
func (b *F32TransferBuf) WriteAt(floats []float32, off int) {
	copy(b.buffer[off:], floats)
}

// Get copy bytes to js and return the reference
func (b *F32TransferBuf) Get() js.Value {
	js.CopyBytesToJS(b.jsBytes, F32Bytes(b.buffer...))
	return b.jsBytes
}

// F32Bytes Cast slice of floats to slice of bytes
func F32Bytes(values ...float32) []byte {
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&values))
	header.Len *= f32size
	header.Cap *= f32size

	// Convert slice header to []byte
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
