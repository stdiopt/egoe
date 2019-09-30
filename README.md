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

## Messaging/callbacks `github.com/stdiopt/egoe/ecs`

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
