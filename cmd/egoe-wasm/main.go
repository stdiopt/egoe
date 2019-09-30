package main

import (
	"log"

	"github.com/stdiopt/egoe"
	"github.com/stdiopt/egoe/ecs"
	"github.com/stdiopt/egoe/systems/dom"
)

func main() {
	log.SetFlags(0)
	m := ecs.New(dom.System, egoe.CustomSystem, egoe.StatSystem)

	m.Start()

	select {}

}
