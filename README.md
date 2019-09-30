# ecs go experiments (egoe)

Exploring `entity component system` pattern in go

After a couple of tries/benchmarks on different methods

## Running

Using hajimehoshi wasm serve https://github.com/hajimehoshi/wasmserve

Get the wasmserve app

```bash
go get -u github.com/hajimehoshi/wasmserve
```

Run on repo path

```
wasmserve ./cmd/egoe
```

## Strong typed components and entities

System is no more no less than a func that accepts a `ecs.Manager` and
registers handlers for specific types on the manager

Component any struct that will be embed on entities. Component structs usually
have methods to return it self so the systems are able to interface and filter

Entity any struct with certain components

## Messaging/callbacks

`github.com/stdiopt/egoe/ecs`

It uses reflection to match specific types and perform the calls,
after benchmarking for a while with different scenarios the results weren't
much different using a reflection Call, although it performs an allocation while
using the reflect.Call internally

```go
m := ecs.New()

m.Watch(func(anystruct){})
m.Handle(func(anystruct){})
m.Query(func(anystruct){})

m.Trigger(anystruct)
```

## Packages

- github.com/stdiopt/egoe/ecs it is the hub where systems interacts with each
  others, a kind of messaging bus
- github.com/stdiopt/egoe/gl contains code from `golang.org/x/mobile/gl` mostly
  the interface and gl constants and a partially implemented webgl wrapper as
  most of things were implemented on demand
- github.com/stdiopt/egoe/systems/dom system that prepares the HTML page with the
  canvas and translate the user inputs
- github.com/stdiopt/egoe/systems/input currently it only have event types
  related to pointer and keyboard
- github.com/stdiopt/egoe/systems/renderer intended to be platform agnostic but
  it currently have some exceptions for buffer transfer (wasm,
  js.CopyBytesToJS..) it will render entities with certain specific components
  world.{Transform,Operation,Material}
- github.com/stdiopt/egoe/world - contains the common data only components and
  events for the world {Transform,Operation,Material}
- github.com/stdiopt/egoe - the domain specific system and entities
