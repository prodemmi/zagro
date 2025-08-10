package zagro

import (
	"errors"
	"sync"
)

// ZagroMessage represents the message passed to event listeners.
// It holds arbitrary data as the event payload.
type ZagroMessage struct {
	Data any `json:"data"`
}

// ZagroCallback defines the function signature for event listeners.
// It receives a pointer to a ZagroMessage.
type ZagroCallback func(*ZagroMessage)

// listener wraps a callback with a unique ID for identification.
type listener struct {
	id       int
	callback ZagroCallback
}

// ZagroOptions allows configuration of a Zagro event emitter instance.
type ZagroOptions struct {
	// MaxListeners limits how many listeners can be registered per event.
	// 0 means unlimited listeners.
	MaxListeners int
}

// Zagro is a concurrency-safe event emitter inspired by JavaScript's EventEmitter.
// It supports multiple listeners per event, one-time listeners, and listener removal.
type Zagro struct {
	mu           sync.Mutex
	events       map[string][]listener
	nextID       int
	maxListeners int
}

// NewZagro creates a new Zagro event emitter.
// Accepts optional ZagroOptions such as MaxListeners to limit listeners per event.
func NewZagro(opts ...ZagroOptions) *Zagro {
	em := &Zagro{
		events: make(map[string][]listener),
	}
	if len(opts) > 0 {
		em.maxListeners = opts[0].MaxListeners
	}
	return em
}

// On registers a new listener callback for the specified event.
// Returns a unique listener ID for later removal and an error if max listeners exceeded.
func (e *Zagro) On(event string, cb ZagroCallback) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.maxListeners > 0 && len(e.events[event]) >= e.maxListeners {
		return 0, errors.New("max listeners exceeded for event: " + event)
	}

	e.nextID++
	l := listener{id: e.nextID, callback: cb}
	e.events[event] = append(e.events[event], l)
	return l.id, nil
}

// Once registers a one-time listener for the event.
// The listener is automatically removed after the first invocation.
func (e *Zagro) Once(event string, cb ZagroCallback) (int, error) {
	var id int
	var err error

	wrapper := func(msg *ZagroMessage) {
		e.Off(event, id)
		cb(msg)
	}

	id, err = e.On(event, wrapper)
	return id, err
}

// Emit triggers all listeners registered for the given event,
// passing the provided ZagroMessage to each callback.
func (e *Zagro) Emit(event string, msg *ZagroMessage) {
	e.mu.Lock()
	callbacks := make([]ZagroCallback, 0, len(e.events[event]))
	for _, l := range e.events[event] {
		callbacks = append(callbacks, l.callback)
	}
	e.mu.Unlock()

	for _, cb := range callbacks {
		cb(msg)
	}
}

// Off removes a specific listener from an event by its unique ID.
func (e *Zagro) Off(event string, id int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	listeners, ok := e.events[event]
	if !ok {
		return
	}
	for i, l := range listeners {
		if l.id == id {
			e.events[event] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
	if len(e.events[event]) == 0 {
		delete(e.events, event)
	}
}

// RemoveAll removes all listeners registered for the specified event.
func (e *Zagro) RemoveAll(event string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.events, event)
}

// Count returns the number of listeners registered for a specific event.
func (e *Zagro) Count(event string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.events[event])
}

// CountAll returns the total number of listeners registered across all events.
func (e *Zagro) CountAll() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	count := 0
	for _, listeners := range e.events {
		count += len(listeners)
	}
	return count
}
