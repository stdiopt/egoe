package ecs

import (
	"reflect"
	"strings"
	"time"
)

// Handler holds the func callback and some extra fields for profiling
type Handler struct {
	Desc      string
	CallStart time.Time
	CallEnd   time.Time
	Fn        interface{}
}

// Describe the handler for debug purposes
func (h *Handler) Describe(v string) *Handler {
	h.Desc = v
	return h
}

// Entry holds the list of handlers for a specific function signature
type Entry struct {
	Type      reflect.Type
	Last      interface{}
	CallStart time.Time
	CallEnd   time.Time
	Handlers  []*Handler
}

// Messaging thing
type Messaging struct {
	Entries map[reflect.Type]*Entry
}

// ReflAuto registers func handelrs by checking Methods on a struct
// If method has Handler prefix it will register as an handler
// If method has Watch prefix it will register as a watcher
func (m *Messaging) ReflAuto(s interface{}) {
	f := reflect.ValueOf(s)
	sn := f.Elem().Type().Name()
	for i := 0; i < f.NumMethod(); i++ {
		mn := f.Type().Method(i).Name
		if strings.HasPrefix(mn, "Handle") {
			m.Handle(f.Method(i).Interface()).Describe(sn + "." + mn)
			continue
		}
		if strings.HasPrefix(mn, "Watch") {
			m.Watch(f.Method(i).Interface()).Describe(sn + "." + mn)
		}
	}
}

// Handle registers an handler for the fn input type and returns the registered
// handler
// the fn interface must be a function with 1 parameter and no returns
// (i.e: func(type))
func (m *Messaging) Handle(fn interface{}) *Handler {
	entry := m.entryFor(fnTyp(fn), true)

	h := &Handler{Fn: fn}
	entry.Handlers = append(entry.Handlers, h)
	return h
}

// Watch registers an handle that triggers upon registration if there is a last
// message, it will return the registered handler
// the fn interface must be a function with 1 parameter and no returns
// (i.e: func(type))
func (m *Messaging) Watch(fn interface{}) *Handler {
	entry := m.entryFor(fnTyp(fn), true)

	h := &Handler{Fn: fn}
	entry.Handlers = append(entry.Handlers, h)

	if entry.Last != nil {
		args := []reflect.Value{reflect.ValueOf(entry.Last)}
		reflect.ValueOf(fn).Call(args)
	}
	return h
}

// Query will callback if the signature type has a message
// the fn interface must be a function with 1 parameter and no returns
// (i.e: func(type))
func (m *Messaging) Query(fn interface{}) {
	entry := m.entryFor(fnTyp(fn), false)

	if entry == nil || entry.Last == nil {
		return
	}

	args := []reflect.Value{reflect.ValueOf(entry.Last)}
	reflect.ValueOf(fn).Call(args)
}

// Trigger will retrieve handlers for the specific signature type and call each
// one
//
// XXX: Might be bad to hold a reference to last value while triggering since
// it can be a huge value, consider creating a new kind of Trigger method for this
// purpose like `Persist(v interface)`
func (m *Messaging) Trigger(v interface{}) {
	k := reflect.TypeOf(v)

	e := m.entryFor(k, true)
	e.Last = v

	e.CallStart = time.Now()

	args := []reflect.Value{reflect.ValueOf(v)}

	for _, h := range e.Handlers {
		h.CallStart = time.Now()
		reflect.ValueOf(h.Fn).Call(args)
		h.CallEnd = time.Now()
	}
	e.CallEnd = time.Now()
}

// entryFor retrieves the entry for the type, if create is true it will create it
func (m *Messaging) entryFor(k reflect.Type, create bool) *Entry {
	if m.Entries == nil {
		if !create {
			return nil
		}
		m.Entries = map[reflect.Type]*Entry{}
	}

	e, ok := m.Entries[k]
	if !ok {
		if !create {
			return nil
		}

		e = &Entry{
			Type:     k,
			Handlers: []*Handler{},
		}
		m.Entries[k] = e
	}

	return e
}

// fnTyp retrieves the dominant type for the specific signature func(type)
func fnTyp(fn interface{}) reflect.Type {
	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func && typ.NumIn() != 1 {
		panic("wrong type, should be a func with 1 param")
	}
	return typ.In(0)
}
