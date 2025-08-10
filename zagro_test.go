package zagro

import (
	"strconv"
	"sync"
	"testing"
)

func TestEventEmitter(t *testing.T) {
	em := NewZagro(ZagroOptions{MaxListeners: 2})

	var mu sync.Mutex
	received := make(map[string]int)

	handler := func(id string) ZagroCallback {
		return func(msg *ZagroMessage) {
			mu.Lock()
			defer mu.Unlock()
			received[id]++
		}
	}

	// تست تابع On و Emit
	id1, err := em.On("event1", handler("h1"))
	if err != nil {
		t.Fatal("On returned error:", err)
	}
	_, err = em.On("event1", handler("h2"))
	if err != nil {
		t.Fatal("On returned error:", err)
	}

	// اضافه کردن listener سوم باید ارور بده (maxListeners=2)
	_, err = em.On("event1", handler("h3"))
	if err == nil {
		t.Error("expected error for maxListeners exceeded, got nil")
	}

	// Emit باید هر دو listener را صدا بزند
	em.Emit("event1", &ZagroMessage{Data: "test"})
	em.Emit("event1", &ZagroMessage{Data: "test"})

	mu.Lock()
	if received["h1"] != 2 || received["h2"] != 2 {
		t.Errorf("unexpected counts: %v", received)
	}
	mu.Unlock()

	// تست Once
	onceCalled := 0
	_, err = em.Once("event2", func(msg *ZagroMessage) {
		onceCalled++
	})
	if err != nil {
		t.Fatal(err)
	}

	em.Emit("event2", &ZagroMessage{})
	em.Emit("event2", &ZagroMessage{})
	if onceCalled != 1 {
		t.Errorf("Once listener called %d times, expected 1", onceCalled)
	}

	// تست Count و CountAll
	if c := em.Count("event1"); c != 2 {
		t.Errorf("Count(event1) = %d, want 2", c)
	}
	if c := em.Count("event2"); c != 0 {
		t.Errorf("Count(event2) = %d, want 0 after once called", c)
	}
	if c := em.CountAll(); c != 2 {
		t.Errorf("CountAll = %d, want 2", c)
	}

	// تست Off
	em.Off("event1", id1)
	if c := em.Count("event1"); c != 1 {
		t.Errorf("Count(event1) after Off = %d, want 1", c)
	}

	// تست RemoveAll
	em.RemoveAll("event1")
	if c := em.Count("event1"); c != 0 {
		t.Errorf("Count(event1) after RemoveAll = %d, want 0", c)
	}
}

func TestMultipleEvents(t *testing.T) {
	em := NewZagro()

	counts := make(map[string]int)
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		eventName := "event" + strconv.Itoa(i)
		_, _ = em.On(eventName, func(e string) ZagroCallback {
			return func(msg *ZagroMessage) {
				mu.Lock()
				counts[e]++
				mu.Unlock()
			}
		}(eventName))
	}

	for i := 0; i < 5; i++ {
		em.Emit("event"+strconv.Itoa(i), &ZagroMessage{Data: i})
	}

	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < 5; i++ {
		eventName := "event" + strconv.Itoa(i)
		if counts[eventName] != 1 {
			t.Errorf("event %s called %d times, want 1", eventName, counts[eventName])
		}
	}
}
