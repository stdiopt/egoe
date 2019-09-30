package egoe

import (
	"log"
	"math"
	"math/rand"

	"github.com/stdiopt/egoe/ecs"
	"github.com/stdiopt/egoe/gl"
	"github.com/stdiopt/egoe/systems/input"
	"github.com/stdiopt/egoe/world"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	nThings = 5000
	minDist = float32(5)
	areaX   = float32(15)
	areaY   = float32(10)
)

// Cool aliases
type (
	vec2 = mgl32.Vec2
	vec3 = mgl32.Vec3
	vec4 = mgl32.Vec4
	mat4 = mgl32.Mat4
)

// Renderable Component
type Renderable struct {
	world.TransformComponent
	world.OperationComponent
	world.MaterialComponent
}

// Thing entity is a mover unit in screen controlled by the customSystem
type Thing struct {
	Renderable
	customComponent
}

// Reset the thing
func (t *Thing) Reset(sz vec2) {
	mat := t.Material()
	custom := t.custom()
	transform := t.Transform()

	if mat.Props == nil {
		mat.Props = map[string]interface{}{}
	}
	mat.Props["color"] = world.Color{}
	custom.speed = rand.Float32() * 0.2
	custom.dir = rand.Float32() * math.Pi * 2
	custom.life = 1
	custom.lifeS = rand.Float32() * 0.01
	transform.World.Translation[0] = rand.Float32()*sz[0]*2 - sz[0]
	transform.World.Translation[1] = rand.Float32()*sz[1]*2 - sz[1]
}

// Custom component
type customComponent struct {
	turner float32
	dir    float32
	speed  float32
	life   float32
	lifeS  float32
}

func (c *customComponent) custom() *customComponent { return c }

// CustomSystem thing
func CustomSystem(m *ecs.Manager) {

	var (
		winSize      vec2
		targets      = []vec2{}
		projection   mat4
		bgShape      *Renderable
		things       []*Thing
		pointerShape *Renderable

		camRotY float32
		camRotX float32
		totalDt float32
	)
	log.Println("Init Custom System")

	// Camera stuff
	m.Watch(func(evt world.ResizeEvent) {
		winSize = vec2(evt)
		cam := mgl32.LookAtV(
			vec3{0, -2, 10}, // position
			vec3{0, -1, 0},  // lookAt
			vec3{0, 1, 0},   // up
		)
		projection = mgl32.Perspective(math.Pi/2, winSize[0]/winSize[1], 0.1, 1000).Mul4(cam)
		m.Trigger(world.CameraEvent(projection))
	})

	// More than one
	m.Handle(func(evt input.PointerEvent) {
		PVInv := projection.Inv()
		targets = []vec2{}
		for i, p := range evt.Pointers {
			ndc := vec4{2*p[0]/winSize[0] - 1, 1 - 2*p[1]/winSize[1], 1, 1}
			dir := PVInv.Mul4x1(ndc).Vec3().Normalize()

			cp := PVInv.Col(3) // Camera position
			nv := intersectPoint(
				dir,
				vec3{cp[0] / cp[3], cp[1] / cp[3], cp[2] / cp[3]},
				vec3{0, 0, 1},
				vec3{0, 0, 0},
			)
			targets = append(targets, nv.Vec2())

			if i == 0 {
				pointerShape.Transform().World.Translation = nv
				pointerShape.Transform().Local.Scale = vec3{minDist, minDist, 1}
			}
		}
	})
	m.Handle(func(evt world.UpdateEvent) {
		dt := float32(evt) * 30 // TimeScale
		totalDt += dt

		projection = projection.
			Mul4(mgl32.HomogRotate3DY(camRotY * dt)).
			Mul4(mgl32.HomogRotate3DX(camRotX * dt))
		m.Trigger(world.CameraEvent(projection))

		// Update camera
		pointerShape.Transform().Local.Rotation = vec3{0, 0, totalDt * 0.01}
		for _, t := range things {
			custom := t.custom()
			transform := t.Transform()
			material := t.Material()

			custom.life -= custom.lifeS

			if custom.life <= 0 {
				t.Reset(vec2{areaX, areaY})
			}
			speed := custom.speed
			nearest := float32(1000)
			dir := custom.dir
			for _, target := range targets {
				// Dist from Point
				dx := target[0] - transform.World.Translation[0]
				dy := target[1] - transform.World.Translation[1]
				dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
				if dist >= nearest {
					continue
				}
				dir = float32(math.Atan2(float64(dy), float64(dx)))
				nearest = dist
			}
			opacity := float32(math.Sin(float64((1 - custom.life) * math.Pi)))
			switch {
			case nearest < 0.1:
				t.Reset(vec2{areaX, areaY})
			case nearest < minDist:
				custom.life = max(0.3, custom.life)
				material.Props["color"] = world.Color{0, 0, 1, opacity}
				custom.dir = dir
				speed = 0.3
			default:
				material.Props["color"] = world.Color{0, 0, 0, opacity}
				custom.turner = limit(custom.turner+(float32(rand.NormFloat64())*0.2), -0.2, 0.2)
				custom.dir += custom.turner * dt
			}

			cos := float32(math.Cos(float64(custom.dir)))
			sin := float32(math.Sin(float64(custom.dir)))

			transform.Local.Rotation[2] = custom.dir
			nx := transform.World.Translation[0] + cos*speed*dt
			ny := transform.World.Translation[1] + sin*speed*dt
			transform.World.Translation[0] = limit(nx, -areaX, areaX)
			transform.World.Translation[1] = limit(ny, -areaY, areaY)

		}
	}).Describe("Update")

	m.Handle(func(evt input.KeyEvent) {
		switch evt.Key {
		case "KeyA":
			camRotY += -0.01
		case "KeyD":
			camRotY += 0.01
		case "KeyW":
			camRotX += 0.01
		case "KeyS":
			camRotX -= 0.01
		}
	})

	m.Handle(func(ecs.StartEvent) {
		// Add entities
		bgShape = makeRenderablePoly("background", 4, gl.TRIANGLE_FAN, vec4{0.5, 1.0, 0.5, 0.8})
		bgShape.TransformComponent.Local.Scale = vec3{areaX + .2, areaY + 0.2, 1}
		m.Entity(bgShape)

		things = createThings(m, vec2{areaX, areaY})
		m.Entity(thingsToEntities(things)...)

		pointerShape = makeRenderablePoly("pointer", 6, gl.TRIANGLE_FAN, vec4{1, 0, 0, 0.2})
		m.Entity(pointerShape)
	})
}

func createThings(m *ecs.Manager, size vec2) []*Thing {
	log.Println("Adding NThings:", nThings)

	// This will be shared between entities
	op := &world.Poly{
		Points: []float32{
			-1, -1, 0,
			1, 0, 0,
			-1, 1, 0,
		},
	}
	ret := []*Thing{}
	// Creating entities
	for i := 0; i < nThings; i++ {
		t := Thing{
			Renderable{
				world.Transform(),
				world.Operation("triangle", op),
				world.MaterialComponent{DrawType: gl.TRIANGLES},
			},
			customComponent{},
		}
		t.TransformComponent.Local.Scale[0] = 0.2
		t.TransformComponent.Local.Scale[1] = 0.2
		t.Reset(size)

		ret = append(ret, &t)
	}
	return ret
}
func thingsToEntities(things []*Thing) []ecs.Entity {
	ret := make([]ecs.Entity, len(things))
	for i := range things {
		ret[i] = things[i]
	}
	return ret
}
