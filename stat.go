package egoe

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"sort"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stdiopt/egoe/ecs"
	"github.com/stdiopt/egoe/systems/dom"
	"github.com/stdiopt/egoe/systems/input"
	"github.com/stdiopt/egoe/world"
)

// StatSystem thing
func StatSystem(m *ecs.Manager) {

	// Dom, wasm specific
	statDom := dom.El("pre",
		dom.Attr{"style": `
			position:fixed;
			padding:10px;
			margin:10px;
			background:rgba(0,0,0,0.8);
			color:white;
			pointer-events:none;
			line-height: 1.5em;
			font-size:0.7em;
		`},
	)
	dom.Body.Call("appendChild", statDom)

	// Profiling
	go func() {
		for {
			statStr := statUpdate(m)

			statDom.Set("innerHTML", statStr)
			log.Println(statStr)

			time.Sleep(time.Second * 3)
		}
	}()

	m.Handle(func(evt input.KeyEvent) {
		if evt.Type == input.KeyUp && evt.Key == "F10" {
			statUpdate(m)
		}
	}).Describe("stat key")

}

// Read stats into a formated string
func statUpdate(m *ecs.Manager) string {
	var winSize mgl32.Vec2
	m.Query(func(evt world.ResizeEvent) {
		winSize = mgl32.Vec2(evt)
	})

	memStat := runtime.MemStats{}
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "Width: %.2f, Height: %.2f\n", winSize[0], winSize[1])
	runtime.ReadMemStats(&memStat)
	fmt.Fprintf(buf,
		"GC: %v Pause: %v, CurMem: %.2fk\n",
		memStat.NumGC,
		time.Duration(memStat.PauseNs[(memStat.NumGC+255)%256]),
		float64(memStat.Alloc)/1024,
	)

	// Since its a map we maintain an order
	entries := []*ecs.Entry{}
	for _, e := range m.Messaging.Entries {
		entries = append(entries, e)
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Type.Name() < entries[j].Type.Name()
	})

	for _, e := range entries {
		dt := e.CallEnd.Sub(e.CallStart)
		fmt.Fprintf(buf, "Handler: %s delta: %v fps: %.2f\n", e.Type.Name(), dt.Round(time.Millisecond/100), float64(time.Second)/float64(dt))
		for _, h := range e.Handlers {
			fmt.Fprintf(buf, "  delta %v: %s\n", h.CallEnd.Sub(h.CallStart), h.Desc)
		}
	}
	return buf.String()
}
