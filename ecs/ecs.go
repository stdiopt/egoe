package ecs

// Entity should be a struct of any kind
type Entity interface{}

// SystemFunc type of function to initialize a system
type SystemFunc func(*Manager)

// Manager the root thing
type Manager struct {
	Messaging
	// XXX: If we want to track entities
	//Entities []Entity
}

// New create a new ECS Manager
func New(fns ...SystemFunc) *Manager {
	m := &Manager{}

	for _, fn := range fns {
		fn(m)
	}
	return m
}

// Start thing
func (m *Manager) Start() {
	m.Trigger(StartEvent{})
}

// Entity adds an entity
func (m *Manager) Entity(e ...Entity) {
	m.Trigger(EntitiesAddEvent(e))
}

// Destroy an entity
func (m *Manager) Destroy(e ...Entity) {
	m.Trigger(EntitiesDestroyEvent(e))
}

// EntitiesAddEvent is triggered when entities are added
type EntitiesAddEvent []Entity

// EntitiesDestroyEvent is triggered when entities are destroyed
type EntitiesDestroyEvent []Entity

// StartEvent fired when things starts
type StartEvent struct{}
