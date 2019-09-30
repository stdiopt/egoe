// +build js,wasm

package dom

import (
	"log"
	"strings"
	"syscall/js"

	"github.com/stdiopt/egoe/ecs"
	"github.com/stdiopt/egoe/gl"
	"github.com/stdiopt/egoe/systems/input"
	"github.com/stdiopt/egoe/systems/renderer"
	"github.com/stdiopt/egoe/world"
)

func init() {
	// Override head
	Document.Get("head").Set("innerHTML", `
	<meta name="mobile-web-app-capable" content="yes">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="manifest" href="/manifest.webmanifest">
	<title> ECS research </title>
	<style>
		* {box-sizing: border-box;}
		body{ height: 100vh; margin:0; padding:0; }
		canvas { position:fixed; top: 0px; height:100%; width: 100%; }
		#fs-btn { z-index:10; position:fixed; top:5px;right:5px; }
	</style>
	`)
}

// System for ecs manager
func System(m *ecs.Manager) {
	s := system{manager: m}
	s.init()
}

type system struct {
	manager *ecs.Manager
	canvas  js.Value

	Width, Height float64
}

func (s *system) init() {

	fullScreenBtn := El("button", Attr{"id": "fs-btn"}, Text("fullscreen"))
	fullScreenBtn.Call("addEventListener", "click", js.FuncOf(func(t js.Value, args []js.Value) interface{} {
		Body.Call("requestFullscreen")
		return nil
	}))
	Body.Call("appendChild", fullScreenBtn)

	s.canvas = El("canvas")
	Body.Call("appendChild", s.canvas)

	// Get gl Context from WebGL thingy
	opt := js.Global().Get("Object").New()
	opt.Set("preserveDrawingBuffer", true)
	opt.Set("antialias", true)
	g := gl.GetWebGL(s.canvas.Call("getContext", "webgl2", opt))

	renderer.System(g, s.manager)
	//r := &glRendererS{}
	//r.Init(s)

	s.checkCanvasSize()
	s.setupEvents()
}

func (s *system) checkCanvasSize() {
	size := s.canvas.Call("getBoundingClientRect")
	w := size.Get("width").Float()
	h := size.Get("height").Float()
	if w != s.Width || h != s.Height {
		s.canvas.Set("width", w)
		s.canvas.Set("height", h)
		s.Width, s.Height = w, h
		s.manager.Trigger(world.ResizeEvent{float32(w), float32(h)})
		log.Printf("Canvas resize: %vx%v", w, h)
	}
}

var (
	evtMap = map[string]input.PointerType{
		"mousedown":   input.MouseDown,
		"mouseup":     input.MouseUp,
		"mousemove":   input.MouseMove,
		"touchstart":  input.PointerDown,
		"touchmove":   input.PointerMove,
		"touchcancel": input.PointerCancel,
		"touchend":    input.PointerEnd,
	}
)

func (s *system) setupEvents() {
	m := s.manager

	log.Println("Registering pointer events")
	ptrEvent := js.FuncOf(s.handlePointerEvent)
	for k := range evtMap {
		s.canvas.Call("addEventListener", k, ptrEvent)
	}

	log.Println("Registering key events")
	keyEvent := js.FuncOf(s.handleKeyEvent)
	js.Global().Call("addEventListener", "keydown", keyEvent)
	js.Global().Call("addEventListener", "keyup", keyEvent)
	js.Global().Call("addEventListener", "keypress", keyEvent)

	var prevFrameTime float64
	var ticker js.Func
	ticker = js.FuncOf(func(t js.Value, args []js.Value) interface{} {
		s.checkCanvasSize()
		dt := args[0].Float()
		dtSec := (dt - prevFrameTime) / 1000
		m.Trigger(world.UpdateEvent(dtSec))
		prevFrameTime = dt
		js.Global().Call("requestAnimationFrame", ticker)
		return nil
	})
	js.Global().Call("requestAnimationFrame", ticker)
}

var (
	keyEvtMap = map[string]int{
		"keydown": input.KeyDown,
		"keyup":   input.KeyUp,
	}
)

func (s *system) handleKeyEvent(t js.Value, args []js.Value) interface{} {
	evt := args[0]
	code := evt.Get("code").String()
	etype := evt.Get("type").String()

	if code == "F12" {
		return nil
	}
	evt.Call("preventDefault")
	s.manager.Trigger(input.KeyEvent{keyEvtMap[etype], code})

	return nil
}
func (s *system) handlePointerEvent(t js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("preventDefault")
	etype := evt.Get("type").String()

	cevt := input.PointerEvent{}

	switch {
	case strings.HasPrefix(etype, "mouse"):
		cevt.Type = evtMap[etype]
		cevt.Pointers = map[int]input.PointerData{
			0: {
				float32(evt.Get("pageX").Float()),
				float32(evt.Get("pageY").Float()),
			},
		}
	case strings.HasPrefix(etype, "touch"):
		cevt.Type = evtMap[etype]
		pts := map[int]input.PointerData{}

		touches := evt.Get("changedTouches")
		for i := 0; i < touches.Length(); i++ {
			t := touches.Index(i)
			id := t.Get("identifier").Int()
			pts[id] = input.PointerData{
				float32(t.Get("pageX").Float()),
				float32(t.Get("pageY").Float()),
			}
		}
		touches = evt.Get("touches")
		for i := 0; i < touches.Length(); i++ {
			t := touches.Index(i)
			id := t.Get("identifier").Int()
			pts[id] = input.PointerData{
				float32(t.Get("pageX").Float()),
				float32(t.Get("pageY").Float()),
			}
		}
		cevt.Pointers = pts
	}

	s.manager.Trigger(cevt)

	return nil
}
