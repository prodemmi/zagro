# Zagro

A lightweight, concurrency-safe **EventEmitter** implementation for Go, inspired by JavaScript’s EventEmitter.

---

## Overview

`zagro` provides a simple and efficient event system in Go, allowing you to register event listeners, emit events, and manage listeners with support for one-time listeners and listener limits — all safely usable in concurrent environments.

---

## Features

- Register listeners for events (`On`)
- Register one-time listeners (`Once`)
- Emit events to all registered listeners (`Emit`)
- Remove specific listeners by ID (`Off`)
- Remove all listeners for an event (`RemoveAll`)
- Count listeners per event (`Count`)
- Count total listeners across all events (`CountAll`)
- Optional limit on the number of listeners per event (`MaxListeners`)

---

## Installation

```bash
go get github.com/prodemmi/zagro
```

## Usage
```go
package main

import (
	"fmt"
	"github.com/prodemmi/zagro"
)

func main() {
	emitter := zagro.NewZagro(zagro.ZagroOptions{MaxListeners: 5})

	id, err := emitter.On("greet", func(msg *zagro.ZagroMessage) {
		fmt.Println("Hello,", msg.Data)
	})
	if err != nil {
		panic(err)
	}

	emitter.Emit("greet", &zagro.ZagroMessage{Data: "World"})

	emitter.Off("greet", id)
}
```

# License
This project is licensed under the MIT License.

# Contributing
Feel free to open issues or submit pull requests to improve Zagro.