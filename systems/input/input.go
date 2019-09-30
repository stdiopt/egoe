package input

import "github.com/go-gl/mathgl/mgl32"

// PointerData common
type PointerData mgl32.Vec2

// PointerType pointer event type
type PointerType int

// Pointer comments
const (
	_ = PointerType(iota)
	MouseDown
	MouseUp
	MouseMove
	PointerDown
	PointerMove
	PointerEnd
	PointerCancel
)

// PointerEvent on canvas
type PointerEvent struct {
	Type     PointerType
	Pointers map[int]PointerData
}

// Consts for key handlers
const (
	KeyDown = iota
	KeyUp
)

// KeyEvent thing
type KeyEvent struct {
	Type int
	Key  string // temp, it should be code
}
