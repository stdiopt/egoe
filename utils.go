package egoe

import (
	"math"

	"github.com/stdiopt/egoe/world"

	"github.com/go-gl/mathgl/mgl32"
)

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func limit(v, mn, mx float32) float32 {
	if v < mn {
		return mn
	} else if v > mx {
		return mx
	}
	return v
}
func intersectPoint(rayVector, rayPoint, planeNormal, planePoint mgl32.Vec3) mgl32.Vec3 {
	diff := rayPoint.Sub(planePoint)
	prod1 := diff.Dot(planeNormal)
	prod2 := rayVector.Dot(planeNormal)
	prod3 := prod1 / prod2
	return rayPoint.Sub(rayVector.Mul(prod3))
}

// makeCircle makes a poly circle comprised on
// x:0, y:0, radius: 1
func makeRenderablePoly(name string, n int, drawType uint, color vec4) *Renderable {
	// build poly verts
	vertex := makePoly(n)
	r := &Renderable{
		world.Transform(),
		world.Operation(name, &world.Poly{Points: vertex}),
		world.MaterialComponent{
			DrawType: drawType,
			Props: map[string]interface{}{
				"color": world.Color(color),
			},
		},
	}

	return r
}
func makePoly(n int) []float32 {
	// Special case
	if n == 4 {
		vertex := []float32{
			-1, -1, 0,
			1, -1, 0,
			1, 1, 0,
			-1, 1, 0,
		}
		return vertex
	}
	points := []float32{}
	p := vec3{0, 1, 0}
	theta := float32(math.Pi) / (float32(n) / 2)
	r := mgl32.HomogRotate2D(theta)
	for i := 0; i < n+1; i++ {
		points = append(points, p[:]...)
		p = r.Mul3x1(p)
	}
	return points
}
