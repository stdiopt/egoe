package world

// Operation returns an operation Component
func Operation(name string, op interface{}) OperationComponent {
	return OperationComponent{Name: name, OP: op}
}

// OperationComponent component
type OperationComponent struct {
	Name string
	OP   interface{}
}

// Operation component
func (c *OperationComponent) Operation() *OperationComponent { return c }

// Poly polygon instructions
type Poly struct {
	Points []float32
}

// Line payload for line operation (not in use)
type Line struct {
	From vec2
	To   vec2
}

// Text thing (not in use)
type Text struct {
	FontSize   int // in PX
	Font       string
	Padding    float64 // in PX
	Background string
	Color      string
	Pos        vec2
	Text       string
}
